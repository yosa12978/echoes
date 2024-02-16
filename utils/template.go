package utils

import (
	"fmt"
	"html/template"
	"io"
	"log"
)

func RenderTemplate(w io.Writer, name string, data interface{}) error {
	templPath := fmt.Sprintf("templates/%s.html", name)
	templ, err := template.ParseFiles(
		templPath,
		"templates/header.html",
		"templates/footer.html",
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return templ.Execute(w, data)
}

func RenderBlock(w io.Writer, name string, data interface{}) error {
	templ := template.Must(template.ParseFiles("templates/blocks.html"))
	return templ.ExecuteTemplate(w, name, data)
}
