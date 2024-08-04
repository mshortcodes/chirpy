package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := r.PathValue("chirpID")
	chirpId, err := strconv.Atoi(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp doesn't exist")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	authorIDStr := r.URL.Query().Get("author_id")
	hasAuthor := false
	if authorIDStr != "" {
		hasAuthor = true
	}

	if hasAuthor {
		authorIDInt, _ := strconv.Atoi(authorIDStr)
		chirps, err = cfg.DB.GetChirpsByAuthor(authorIDInt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
			return
		}
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	sortBy := r.URL.Query().Get("sort")

	if sortBy == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
