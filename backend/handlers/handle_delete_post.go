package handlers

import (
	"blogsite/internal/database"
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// DeletePost handles DELETE /api/posts/{id}
// Only the author of the post can delete it.
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	existing, err := h.Queries.GetPost(r.Context(), postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existing.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.Queries.DeletePost(r.Context(), database.DeletePostParams{
		ID:     postID,
		UserID: userID,
	}); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
