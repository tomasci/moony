package queries_client

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"moony/database/sqlc"
	"net/url"
	"os"
	"sync"
)

var (
	// sqlc queries instance
	queries *sqlc.Queries
	// db connection pool instance
	pool *pgxpool.Pool
	// once - to ensure instance created only once
	once sync.Once
	// for errors
	initError error
)

func getConnectionURL() *url.URL {
	// get keys from env
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresDb := os.Getenv("POSTGRES_DB")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")

	// create db connection url
	dsn := url.URL{
		Scheme: "postgres",
		Host:   postgresHost,
		Path:   postgresDb,
		User:   url.UserPassword(postgresUser, postgresPassword),
	}

	return &dsn
}

func InitializeDatabase() {
	once.Do(func() {
		log.Println("Initializing database connection pool")
		ctx := context.Background()
		connectionUrl := getConnectionURL()

		// pool config
		config, err := pgxpool.ParseConfig(connectionUrl.String())
		if err != nil {
			log.Fatalf("failed to parse connection url: %v", err)
		}

		pool, initError = pgxpool.NewWithConfig(ctx, config)
		if initError != nil {
			log.Fatalf("failed to create database connection pool: %v", initError)
		}

		// use sqlc
		queries = sqlc.New(pool)
	})
}

// GetDBConnectionPool – global access to pgx database connection
func GetDBConnectionPool() (*pgxpool.Pool, error) {
	if pool == nil {
		InitializeDatabase()
	}

	return pool, initError
}

// GetQueriesClient – global access to sqlc.Queries intance
func GetQueriesClient() (*sqlc.Queries, error) {
	if pool == nil {
		InitializeDatabase()
	}

	return queries, initError
}
