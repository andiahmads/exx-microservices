package main

import (
	"authentication-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	//connect to database
	conn := connectTODB()
	if conn == nil {
		log.Panic("Cant connect to database")
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	fmt.Println("Starting authentication service")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectTODB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready")
			counts++
		} else {
			log.Println("Connect to Postgres")
			return connection

		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("backing off for two second...")
		time.Sleep(2 * time.Second)
		continue

	}
}
