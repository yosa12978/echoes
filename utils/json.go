package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/validation"
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

func ReadJsonAndValidate[T validation.Validatable[T]](ctx context.Context, body io.Reader) (T, map[string]string, error) {
	var dest T
	if err := json.NewDecoder(body).Decode(&dest); err != nil {
		return dest, nil, err
	}
	res, problems, ok := dest.Validate(ctx)
	if !ok {
		return dest, problems,
			types.NewErrValidationFailed(errors.New("validation failed"))
	}
	return res, nil, nil
}
