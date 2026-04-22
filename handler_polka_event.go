package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/janeturovskaya/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaEvent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "api key error", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, 401, "wrong api key", errors.New("Wrong API key"))
		return
	}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can not decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "Error parsing user ID", err)
		return
	}

	_, err = cfg.db.UpgradeSubscription(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "user not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "db error", err)
		return
	}

	respondWithJSON(w, 204, nil)

}
