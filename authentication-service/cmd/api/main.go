package main

import (
	"authentication/cmd/api/routes"
	"authentication/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

func main() {
	log.Println("Starting the application")

	conn := connectToDB()
	if conn == nil {
		log.Fatal("Could not connect to the database")
	}

	app := routes.Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Error connecting to the database", err)
			counts++
		} else {
			log.Println("Connected to the database")
			return connection
		}

		if counts > 10 {
			log.Println("Could not connect to the database after 10 attempts")
			return nil
		}

		log.Println("Retrying connection to the database")
		time.Sleep(2 * time.Second)
		continue
	}
}
