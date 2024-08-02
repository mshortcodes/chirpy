package main

import (
	"chirpy/internal/auth"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpIDInt, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	userIDStr, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse User ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpIDInt)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	if userIDInt != chirp.AuthorID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp!")
		return
	}

	err = cfg.DB.DeleteChirp(chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
