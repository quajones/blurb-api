package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB() (DB, error) {
	connStr := getConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return DB{}, err
	}
	if err = db.Ping(); err != nil {
		return DB{}, err
	}
	return DB{db}, nil
}

func getConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"), getEnvOrDefault("DB_PORT", 5432).(int),
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
}

func getEnvOrDefault(variable string, def interface{}) interface{} {
	value, exists := os.LookupEnv(variable)
	if !exists {
		return def
	}
	return value
}
