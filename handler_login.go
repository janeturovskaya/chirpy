package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/janeturovskaya/chirpy/internal/auth"
	"github.com/janeturovskaya/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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

	user, err := apiCfg.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, apiCfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token creating error", err)
	}

	tokenString := auth.MakeRefreshToken()
	refreshToken, err := apiCfg.db.CreateToken(context.Background(), database.CreateTokenParams{
		Token:  tokenString,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}
	type loginResponse struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	respondWithJSON(w, 200, loginResponse{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken.Token,
	})
}
