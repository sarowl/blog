// backend\handlers\handle_signup.go [renamed]
package handlers

import (
	"blogsite/internal/database"
	"database/sql"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	Queries *database.Queries
}

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// This 'user_id' comes from your AuthMiddleware
	firebaseUID := r.Context().Value("user_id").(string)

	var req struct {
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	displayName := sql.NullString{
		String: req.DisplayName,
		Valid:  req.DisplayName != "", // Valid is true if the string is not empty
	}

	user, err := h.Queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:          firebaseUID,
		Email:       req.Email,
		DisplayName: displayName,
	})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toUserResponse(user))
}

func toUserResponse(u database.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName.String, // empty string if NULL/not valid
	}
}
