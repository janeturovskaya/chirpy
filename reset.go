package main

import (
	"context"
	"errors"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden", errors.New("Forbidden error"))
		return
	}

	cfg.fileserverHits.Store(0)
	cfg.db.DeleteUsers(context.Background())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
