package main

import (
	"context"
	"net/http"

	"errors"

	"github.com/google/uuid"
	"github.com/janeturovskaya/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, 404, "Error parsing chirp ID", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, 401, "Token err", err)
		return
	}

	chirp, err := cfg.db.GetChirpByChirpId(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, 404, "No such chirp in db", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, 403, "User is not the author of the chirp", errors.New("User is not the author of the chirp"))
		return
	}

	err = cfg.db.DeleteChirpsByID(context.Background(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No such chirp in db", err)
		return
	}
	respondWithJSON(w, 204, nil)

}
