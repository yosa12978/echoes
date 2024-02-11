package utils

import (
	"fmt"
	"html/template"
	"io"
	"log"
)

func RenderTemplate(w io.Writer, name string, data interface{}) error {
	templPath := fmt.Sprintf("/templates/%s.html", name)
	templ, err := template.ParseFiles("/templates/head.html", templPath)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return templ.Execute(w, data)
}

func RenderBlock(w io.Writer, name string, data interface{}) {
	templ := template.Must(template.ParseFiles("/templates/blocks.html"))
	templ.ExecuteTemplate(w, name, data)
}
