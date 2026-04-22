package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func clearProfineWords(post string) (string, error) {
	var b strings.Builder
	profaneWords := [3]string{"kerfuffle", "sharbert", "fornax"}
	splitted := strings.Split(post, " ")
	for i, w := range splitted {
		lowered := strings.ToLower(w)
		if lowered == profaneWords[0] || lowered == profaneWords[1] || lowered == profaneWords[2] {
			splitted[i] = "****"
		}
		_, err := b.WriteString(splitted[i] + " ")
		if err != nil {
			log.Printf("Error writing string: %v", err)
			return "", err
		}
	}
	res := b.String()
	res = strings.TrimSpace(res)
	return res, nil
}
