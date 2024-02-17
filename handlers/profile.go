package handlers

import (
	"context"
	"net/http"

	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/utils"
)

type Profile interface {
	Get(ctx context.Context) http.Handler
}

type profile struct {
	profileRepo repos.Profile
}

func NewProfile(profileRepo repos.Profile) Profile {
	h := new(profile)
	h.profileRepo = profileRepo
	return h
}

func (h *profile) Get(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := h.profileRepo.Get(ctx)
		utils.RenderBlock(w, "profile", res)
	})
}
