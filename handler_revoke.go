package main

import (
	"context"
	"net/http"

	"github.com/janeturovskaya/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't refresh token", err)
		return
	}

	_, err = cfg.db.RevokeToken(context.Background(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
