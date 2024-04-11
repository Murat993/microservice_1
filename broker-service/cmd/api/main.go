package main

import (
	"broker/cmd/api/routes"
	"log"
	"net/http"
)

const webPort = "80"

func main() {
	app := routes.Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.Routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
