package app

import (
	"context"
	"net/http"

	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/router"
	"github.com/yosa12978/echoes/services"
)

func newServer(ctx context.Context, addr string, logger logging.Logger) http.Server {
	postRepo := repos.NewPostPostgres()
	linkRepo := repos.NewLinkPostgres()
	commentRepo := repos.NewCommentPostgres()
	accountRepo := repos.NewAccountPostgres()
	profileRepo := repos.NewProfileFromConfig()
	announceRepo := cache.NewAnnounceRedis(data.Redis(ctx))

	postService := services.NewPost(
		postRepo,
		cache.NewPostRedis(data.Redis(ctx), logger),
		logger,
		repos.NewPostSearcherPostgres(),
	)
	linkService := services.NewLink(
		linkRepo,
		cache.NewLinkRedis(data.Redis(ctx), logger),
		logger,
	)
	commentService := services.NewComment(
		commentRepo,
		postService,
		cache.NewCommentRedis(data.Redis(ctx), logger),
		logger,
	)
	announceService := services.NewAnnounce(
		announceRepo,
		logger,
	)
	healthService := services.NewHealthService(
		logger,
		data.NewPgPinger(),
		data.NewRedisPinger(ctx),
	)
	accountService := services.NewAccount(accountRepo)
	profileService := services.NewProfile(profileRepo)
	feedService := services.NewFeedService(postService)

	accountService.Seed(ctx)

	router := router.New(
		router.WithLogger(logger),
		router.WithAccountService(accountService),
		router.WithAnnounceService(announceService),
		router.WithCommentService(commentService),
		router.WithFeedService(feedService),
		router.WithLinkService(linkService),
		router.WithPostService(postService),
		router.WithProfileService(profileService),
		router.WithHealthService(healthService),
	)

	return http.Server{
		Addr:    addr,
		Handler: router,
	}
}
