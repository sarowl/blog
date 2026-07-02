package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// This 'user_id' comes from your AuthMiddleware
	firebaseUID := r.Context().Value("user_id").(string)

	user, err := h.Queries.GetUser(r.Context(), firebaseUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toUserResponse(user))
}
