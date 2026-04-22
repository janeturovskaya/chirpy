package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, 404, "Error parsing chirp ID", err)
		return
	}
	chirp, err := apiCfg.db.GetChirpByChirpId(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Error getting chirp by chirp ID", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
