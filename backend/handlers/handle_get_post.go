// backend\handlers\handle_get_post.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// GetPost handles GET /api/posts/{id}
// Public posts are visible to any authenticated user; private posts only to their author.
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idParam := chi.URLParam(r, "id")
	postID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	post, err := h.Queries.GetPost(r.Context(), postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if post.Status != "public" && post.UserID != userID {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// ListPosts handles GET /posts and GET /api/posts
// Public posts are visible to anyone; authenticated requests can also request only the current user's posts via ?mine=true.
// ListPosts handles GET /api/posts (public feed) and /api/posts?mine=true (own posts).
// Requires authentication. The feed itself only returns status='public' posts from any user.
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type postResponse struct {
		ID         uuid.UUID `json:"id"`
		UserID     string    `json:"user_id"`
		AuthorName string    `json:"author_name"`
		Title      string    `json:"title"`
		Body       string    `json:"body"`
		Status     string    `json:"status"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	var posts []postResponse
	if r.URL.Query().Get("mine") == "true" {
		rows, err := h.Queries.ListUserPosts(r.Context(), userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		posts = make([]postResponse, 0, len(rows))
		for _, row := range rows {
			posts = append(posts, postResponse{
				ID:         row.ID,
				UserID:     row.UserID,
				AuthorName: "",
				Title:      row.Title,
				Body:       row.Body,
				Status:     row.Status,
				CreatedAt:  row.CreatedAt,
				UpdatedAt:  row.UpdatedAt,
			})
		}
	} else {
		rows, err := h.Queries.ListPosts(r.Context())
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		posts = make([]postResponse, 0, len(rows))
		for _, row := range rows {
			authorName := row.AuthorName.String
			if authorName == "" {
				authorName = row.UserID
			}
			posts = append(posts, postResponse{
				ID:         row.ID,
				UserID:     row.UserID,
				AuthorName: authorName,
				Title:      row.Title,
				Body:       row.Body,
				Status:     row.Status,
				CreatedAt:  row.CreatedAt,
				UpdatedAt:  row.UpdatedAt,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// ListMyPosts handles GET /api/me/posts and returns only the authenticated user's posts.
func (h *PostHandler) ListMyPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	posts, err := h.Queries.ListUserPosts(r.Context(), userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
