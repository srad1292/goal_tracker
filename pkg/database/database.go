package database

import (
	"fmt"
	"os"
)

func ConnectToDatabase() {
	port := os.Getenv("DB_PORT")
	schema := os.Getenv("DB_SCHEMA")

	fmt.Printf("env: Port = %s, Schema = %s\n", port, schema)
}
