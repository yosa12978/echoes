package repos

import (
	"context"
	"encoding/json"
	"os"

	"github.com/yosa12978/echoes/types"
)

type Profile interface {
	Get(ctx context.Context) types.Profile
}

type profileJson struct {
	profile types.Profile
}

func NewProfileJson(filename string) (Profile, error) {
	repo := new(profileJson)
	file, err := os.Open(filename)
	if err != nil {
		return repo, err
	}
	err = json.NewDecoder(file).Decode(&repo.profile)
	return repo, err
}

func (repo *profileJson) Get(ctx context.Context) types.Profile {
	return repo.profile
}
