package main

import (
	"encoding/json"
	"net/http"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	//check lenght of the data's body
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}
	// everything ok, send respond
	respondWithJSON(w, http.StatusOK, responseBody{
		Valid: true,
	})
}
