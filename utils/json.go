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

func ReadJson(body io.Reader, dest any) error {
	return json.NewDecoder(body).Decode(dest)
}
