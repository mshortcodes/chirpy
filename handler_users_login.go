package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get user")
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:    user.ID,
		Email: user.Email,
	})
}
