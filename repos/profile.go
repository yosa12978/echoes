package repos

import (
	"context"
	"encoding/json"
	"os"

	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/types"
)

type Profile interface {
	Get(ctx context.Context) (*types.Profile, error)
	Update(ctx context.Context, profile types.Profile) (*types.Profile, error)
}

type profileJson struct {
	filename string
	profile  *types.Profile
}

func NewProfileJson(filename string) (Profile, error) {
	repo := new(profileJson)
	file, err := os.Open(filename)
	if err != nil {
		return repo, err
	}
	defer file.Close()
	repo.filename = filename
	err = json.NewDecoder(file).Decode(&repo.profile)
	return repo, err
}

func (repo *profileJson) Get(ctx context.Context) (*types.Profile, error) {
	return repo.profile, nil
}

func (repo *profileJson) Update(ctx context.Context, profile types.Profile) (*types.Profile, error) {
	file, err := os.Open(repo.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	repo.profile = &profile
	return repo.profile, json.NewEncoder(file).Encode(profile)
}

type profileRedis struct {
	cache cache.Cache
}

func NewProfileRedis(cache cache.Cache) Profile {
	return &profileRedis{cache: cache}
}

func (p *profileRedis) Get(ctx context.Context) (*types.Profile, error) {
	return nil, nil
}

func (p *profileRedis) Update(ctx context.Context, profile types.Profile) (*types.Profile, error) {
	return nil, nil
}
