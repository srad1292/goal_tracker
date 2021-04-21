package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbConnection *pgxpool.Pool

func connectToDatabase() *pgxpool.Pool {
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	// db_schema := os.Getenv("DB_SCHEMA")
	db_conn := os.Getenv("DB_MAX_CONNECTIONS")

	var dbURL string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%s", db_user, db_password, db_host, db_port, db_name, db_conn)
	// url = postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = dbpool.Ping(context.Background())

	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	return dbpool

}

func GetDatabase() *pgxpool.Pool {
	if dbConnection == nil {
		dbConnection = connectToDatabase()
	}

	return dbConnection
}
