package router

import (
	"os"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
)

type optionFunc func(*options)

type options struct {
	accountService  services.Account
	announceService services.Announce
	commentService  services.Comment
	feedService     services.Feed
	healthService   services.HealthService
	postService     services.Post
	profileService  services.Profile
	linkService     services.Link
	logger          logging.Logger
}

func defaultOptions() options {
	return options{
		logger: logging.NewJsonLogger(os.Stdout),
	}
}

func newOptions(opts ...optionFunc) options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithLogger(logger logging.Logger) optionFunc {
	return func(o *options) {
		o.logger = logger
	}
}

func WithAccountService(s services.Account) optionFunc {
	return func(o *options) {
		o.accountService = s
	}
}

func WithProfileService(s services.Profile) optionFunc {
	return func(o *options) {
		o.profileService = s
	}
}

func WithHealthService(s services.HealthService) optionFunc {
	return func(o *options) {
		o.healthService = s
	}
}

func WithFeedService(s services.Feed) optionFunc {
	return func(o *options) {
		o.feedService = s
	}
}

func WithAnnounceService(s services.Announce) optionFunc {
	return func(o *options) {
		o.announceService = s
	}
}

func WithCommentService(s services.Comment) optionFunc {
	return func(o *options) {
		o.commentService = s
	}
}

func WithLinkService(s services.Link) optionFunc {
	return func(o *options) {
		o.linkService = s
	}
}

func WithPostService(s services.Post) optionFunc {
	return func(o *options) {
		o.postService = s
	}
}
