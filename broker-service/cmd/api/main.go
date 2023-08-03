package main

import (
	"fmt"
	"log"
	"net/http"
)

type App struct {
	Port string
}

const (
	WEB_PORT = "80"
)

func main() {
	app := App{
		Port: WEB_PORT,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.Port),
		Handler: app.Routes(),
	}
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
