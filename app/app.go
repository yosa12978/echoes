package app

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/config"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/router"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/session"
)

func Run(ctx context.Context) error {
	logger := logging.NewJsonLogger(os.Stdout)
	conn := data.Postgres()
	defer conn.Close()

	rdb := data.Redis(ctx)
	defer rdb.Close()

	session.SetupStore()

	cfg := config.Get()

	server := newServer(
		ctx,
		cfg.Server.Addr,
		logger,
	)

	errch := make(chan error, 1)
	go func() {
		logger.Info("server listening", "addr", cfg.Server.Addr)
		if err := server.ListenAndServe(); err != nil {
			errch <- err
		}
		close(errch)
	}()

	var err error
	select {
	case err = <-errch:
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err = server.Shutdown(timeout)
	}
	return err
}

func newServer(ctx context.Context, addr string, logger logging.Logger) http.Server {
	postRepo := repos.NewPostPostgres()
	linkRepo := repos.NewLinkPostgres()
	commentRepo := repos.NewCommentPostgres()
	accountRepo := repos.NewAccountPostgres()
	profileRepo := repos.NewProfileFromConfig()
	announceRepo := repos.NewAnnounceCache(cache.NewRedisCache(ctx))

	postService := services.NewPost(
		postRepo,
		cache.NewRedisCache(ctx),
		logger,
		repos.NewPostSearcherPostgres(),
	)
	linkService := services.NewLink(
		linkRepo,
		cache.NewRedisCache(ctx),
		logger,
	)
	commentService := services.NewComment(
		commentRepo,
		postService,
		cache.NewRedisCache(ctx),
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
