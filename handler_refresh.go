package main

import (
	"context"

	"net/http"

	"github.com/janeturovskaya/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}

	token, err := cfg.db.GetRefreshTokenByTokenString(context.Background(), tokenString)
	if err != nil {
		respondWithError(w, 401, "Token err", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(context.Background(), token.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user by token from db", err)
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token creating error", err)
		return
	}
	type t struct {
		Token string `json:"token"`
	}
	resp := t{accessToken}

	respondWithJSON(w, 200, resp)
}
