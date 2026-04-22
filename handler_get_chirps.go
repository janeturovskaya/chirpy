package main

import (
	"context"
	"net/http"

	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/janeturovskaya/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	s := r.URL.Query().Get("author_id")
	data := make([]database.Chirp, 0)
	err := errors.New("New error")

	if s == "" {
		data, err = cfg.db.GetChirps(context.Background())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "GetChirps Error", err)
			return
		}
	} else {
		userID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, 404, "Error parsing user ID", err)
			return
		}
		data, err = cfg.db.GetChirpsByAuthor(context.Background(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "GetChirpsByAuthor Error", err)
			return
		}
	}

	chirps := make([]Chirp, 0, len(data))
	for _, chirp := range data {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	sortParam := r.URL.Query().Get("sort")
	switch sortParam {
	case "asc":
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	case "desc":
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	default:
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
