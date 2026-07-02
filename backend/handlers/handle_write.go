package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"blogsite/internal/database"
)

type PostHandler struct {
	Queries *database.Queries
}

type createPostRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status string `json:"status"`
}

type updatePostRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status string `json:"status"`
}

func validateStatus(status string) (string, bool) {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		status = "draft"
	}
	if status != "draft" && status != "public" && status != "private" {
		return "", false
	}
	return status, true
}

// CreatePost handles POST /api/posts
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Body = strings.TrimSpace(req.Body)
	if req.Title == "" || req.Body == "" {
		http.Error(w, "title and body are required", http.StatusBadRequest)
		return
	}

	status, ok := validateStatus(req.Status)
	if !ok {
		http.Error(w, "status must be 'public' or 'private'", http.StatusBadRequest)
		return
	}

	post, err := h.Queries.CreatePost(r.Context(), database.CreatePostParams{
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body,
		Status: status,
	})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// DeletePost handles DELETE /api/posts/{id}
