package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var (
	MaxOpenConnections    = 25
	MaxIdleConnections    = 25
	ConnectionMaxLifetime = 5 * time.Minute
	maxRetries            = 5
	retryDelay            = 5 * time.Second
	connTimeout           = 5 * time.Second
)

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func NewPostgresConnection(dbConfig DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode, int(connTimeout.Seconds()),
	)

	var db *sql.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Attempt %d: Error connecting to PostgreSQL: %s\n", i, err)
			time.Sleep(retryDelay)
			continue
		}

		if err := db.Ping(); err != nil {
			log.Printf("Attempt %d: PostgreSQL ping error: %s\n", i, err)
			time.Sleep(retryDelay)
			continue
		}

		db.SetMaxOpenConns(MaxOpenConnections)
		db.SetMaxIdleConns(MaxIdleConnections)
		db.SetConnMaxLifetime(ConnectionMaxLifetime)

		log.Println("The connection to PostgreSQL has been successfully established")
		return db, nil
	}
	return nil, fmt.Errorf("couldn't connect to the database: %w", err)
}
