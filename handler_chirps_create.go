package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
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

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, userIDInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:       chirp.ID,
		AuthorID: chirp.AuthorID,
		Body:     chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleaned := getCleanedBody(body)
	return cleaned, nil
}

func getCleanedBody(body string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	stars := "****"

	words := strings.Split(body, " ")
	for i, word := range words {
		lowered := strings.ToLower(word)
		_, ok := badWords[lowered]
		if ok {
			words[i] = stars
		}
	}

	cleaned := strings.Join(words, " ")
	return cleaned
}
