package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func WriteJson(w http.ResponseWriter, payload any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(payload)
}

func ReadJson[T any](body io.Reader) (T, error) {
	var dest T
	err := json.NewDecoder(body).Decode(&dest)
	return dest, err
}
