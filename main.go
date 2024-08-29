package main

import (
	"github.com/yosa12978/echoes/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
