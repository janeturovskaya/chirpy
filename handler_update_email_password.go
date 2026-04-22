package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/janeturovskaya/chirpy/internal/auth"
	"github.com/janeturovskaya/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateEmailPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
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

	user, err := cfg.db.GetUserFromRefreshToken(context.Background(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user by token from db", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	newUserData, err := cfg.db.ChangeEmailPassword(context.Background(), database.ChangeEmailPasswordParams{
		Email:    params.Email,
		Password: hash,
		ID:       user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't change email, name in db", err)
		return
	}

	respondWithJSON(w, 200, User{
		ID:          newUserData.ID,
		CreatedAt:   newUserData.CreatedAt,
		UpdatedAt:   newUserData.UpdatedAt,
		Email:       newUserData.Email,
		IsChirpyRed: newUserData.IsChirpyRed,
	})

}
