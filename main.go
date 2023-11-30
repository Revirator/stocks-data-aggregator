package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.ExpandEnv("postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}")
	db, err := DatabaseInit(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	hostAndPort := os.ExpandEnv("${SERVER_HOST}:${SERVER_PORT}")
	server := ServerInit(hostAndPort, db)
	server.Run()
}
