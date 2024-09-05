package utils

import (
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/yosa12978/echoes/config"
	"github.com/yosa12978/echoes/types"
)

var cfg = config.Get()

func RenderView(w io.Writer, view string, title string, payload any) error {
	templPath := fmt.Sprintf("templates/views/%s.html", view)
	templ, err := template.ParseFiles(
		templPath,
		"templates/top.html",
		"templates/bottom.html",
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if title != "" {
		title = "/" + title
	}
	data := types.Templ{
		Title:   cfg.Website.Title + title,
		Logo:    cfg.Website.Logo,
		BgImg:   cfg.Website.BgImg,
		Payload: payload,
	}
	return templ.Execute(w, data)
}

func RenderBlock(w io.Writer, name string, payload any) error {
	templ := template.Must(
		template.ParseFiles(
			"templates/blocks/posts.html",
			"templates/blocks/links.html",
			"templates/blocks/profile.html",
			"templates/blocks/announce.html",
			"templates/blocks/alert.html",
			"templates/blocks/comments.html",
		),
	)
	return templ.ExecuteTemplate(w, name, payload)
}
