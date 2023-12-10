package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/revirator/cfd/companydb"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.ExpandEnv("postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}")
	db := databaseInit(connectionString)

	hostAndPort := os.ExpandEnv("${SERVER_HOST}:${SERVER_PORT}")
	server := ServerInit(hostAndPort, companydb.NewCompanyDatabse(db))
	server.Run()
}

func databaseInit(connectionString string) *sql.DB {
	if db, err := sql.Open("postgres", connectionString); err != nil {
		panic(err)
	} else {
		return db
	}
}
