package main

import (
	"fmt"
	"log"
	"mailhog/cmd/api/routes"
	"net/http"
	"os"
	"strconv"
)

const webPort = "80"

func main() {
	app := routes.Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail server on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Println("Error starting mail server:", err)
	}
}

func createMail() routes.Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := routes.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
	}

	return m
}
