package handlers

import (
	"net/http"

	"github.com/yosa12978/echoes/services"
)

type Feed interface {
	GetFeed() http.HandlerFunc
}

type feed struct {
	feedService services.Feed
}

func NewFeedHandler(feedService services.Feed) Feed {
	return &feed{
		feedService: feedService,
	}
}

func (f *feed) GetFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		feed, err := f.feedService.GenerateFeed(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(feed))
	}
}
