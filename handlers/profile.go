package handlers

import (
	"context"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

type Profile interface {
	Get(ctx context.Context) http.Handler
}

type profile struct {
	profileService services.Profile
	logger         logging.Logger
}

func NewProfile(profileService services.Profile, logger logging.Logger) Profile {
	h := new(profile)
	h.profileService = profileService
	h.logger = logger
	return h
}

func (h *profile) Get(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := h.profileService.Get(ctx)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "can't render profile information")
			return
		}
		utils.RenderBlock(w, "profile", res)
	})
}
