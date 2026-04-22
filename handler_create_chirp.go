package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/janeturovskaya/chirpy/internal/auth"
	"github.com/janeturovskaya/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can not decode parameters", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not authorize", err)
		return
	}
	userID, err := auth.ValidateJWT(tokenString, apiCfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not authorize", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanBody := getCleanedBody(params.Body, badWords)
	chirp, err := apiCfg.db.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   cleanBody,
		UserID: userID,
	})
	if err != nil {
		log.Fatalf("Error creating chirp: %v", err)
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
