package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	polkaApiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing API key")
		return
	}

	if polkaApiKey != cfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Wrong event")
		return
	}

	err = cfg.DB.UpgradeChiryRed(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
