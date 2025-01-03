package queries_client

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"moony/database/sqlc"
	"net/url"
	"os"
	"sync"
)

var (
	// sqlc queries instance
	queries *sqlc.Queries
	// db connection instance
	conn *pgx.Conn
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
		log.Println("Initializing database connection")
		ctx := context.Background()
		connectionUrl := getConnectionURL()
		// try to connect to db
		conn, initError = pgx.Connect(ctx, connectionUrl.String())

		if initError != nil {
			log.Fatalf("error connecting to database: %s", initError)
		}

		// use sqlc
		queries = sqlc.New(conn)
	})
}

// GetDBConnection – global access to pgx database connection
func GetDBConnection() (*pgx.Conn, error) {
	if conn == nil {
		InitializeDatabase()
	}

	return conn, initError
}

// GetQueriesClient – global access to sqlc.Queries intance
func GetQueriesClient() (*sqlc.Queries, error) {
	if conn == nil {
		InitializeDatabase()
	}

	return queries, initError
}
