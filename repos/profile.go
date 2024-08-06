package repos

import (
	"context"
	"encoding/json"
	"os"

	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/config"
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
	cfg   config.Config
}

func NewProfileRedis(cache cache.Cache) Profile {
	return &profileRedis{cache: cache}
}

func (p *profileRedis) Get(ctx context.Context) (*types.Profile, error) {
	profileJson, err := p.cache.Get(ctx, "profile")
	if err != nil {
		profile := types.Profile{
			Name: p.cfg.Profile.Name,
			Bio:  p.cfg.Profile.Bio,
			Icon: p.cfg.Profile.Picture,
		}
		j, _ := json.Marshal(profile)
		_, err := p.cache.Set(ctx, "profile", string(j), -1)
		return &profile, err
	}
	profile := types.Profile{}
	err = json.Unmarshal([]byte(profileJson), &profile)
	return &profile, err
}

func (p *profileRedis) Update(ctx context.Context, profile types.Profile) (*types.Profile, error) {
	return nil, nil
}

type profileFromConfig struct {
	cfg config.Config
}

func NewProfileFromConfig() Profile {
	return &profileFromConfig{
		cfg: config.Get(),
	}
}

// better load profile info from config to redis and change or smth
func (p *profileFromConfig) Get(ctx context.Context) (*types.Profile, error) {
	return &types.Profile{
		Name: p.cfg.Profile.Name,
		Bio:  p.cfg.Profile.Bio,
		Icon: p.cfg.Profile.Picture,
	}, nil
}

func (p *profileFromConfig) Update(ctx context.Context, profile types.Profile) (*types.Profile, error) {
	panic("unimplemented")
}
