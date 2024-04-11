package main

import (
	"log"
	"net/http"
)

const webPort = "8080"

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("Starting broker server on port %s\n", webPort)

	//define http server
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
