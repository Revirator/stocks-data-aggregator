package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hostAndPort := os.ExpandEnv("${SERVER_HOST}:${SERVER_PORT}")
	db := databaseInit()
	server := ServerInit(hostAndPort, db)
	server.Run()
}

func databaseInit() *sql.DB {
	connectionString := os.ExpandEnv("postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}")
	if db, err := sql.Open("postgres", connectionString); err != nil {
		panic(err)
	} else {
		return db
	}
}
