package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("No apiKey found")
	}
	if !strings.HasPrefix(apiKey, "ApiKey ") {
		return "", errors.New("No apiKey found")
	}
	apiKey = strings.TrimPrefix(apiKey, "ApiKey")

	apiKey = strings.Trim(apiKey, " ")
	return apiKey, nil
}
