package services

import (
	"context"

	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Profile interface {
	Get(ctx context.Context) (*types.Profile, error)
}

type profile struct {
	profileRepo repos.Profile
}

func NewProfile(profileRepo repos.Profile) Profile {
	return &profile{profileRepo: profileRepo}
}

func (s *profile) Get(ctx context.Context) (*types.Profile, error) {
	return s.profileRepo.Get(ctx)
}
