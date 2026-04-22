package auth

import (
	"errors"
	"net/http"
	"strings"
)

//This function should look for the Authorization header in the headers parameter and return the TOKEN_STRING
//  if it exists (stripping off the Bearer prefix and whitespace).
// If the header doesn't exist, return an error.
// This is an easy one to write a unit test for, and I'd recommend doing so.

func GetBearerToken(headers http.Header) (string, error) {

	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("No token found")
	}
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer")
	}

	tokenString = strings.Trim(tokenString, " ")
	return tokenString, nil
}
