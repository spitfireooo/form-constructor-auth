package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type DBConfig struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
	SSLMode  string
}

const (
	UsersTable  = "users"
	TokensTable = "tokens"
)

var Connect *sqlx.DB

func DatabaseInit(con DBConfig) error {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		con.Username, con.Password, con.Database, con.Host, con.Port, con.SSLMode,
	)

	if db, err := sqlx.Connect("postgres", dsn); err != nil {
		return err
	} else {
		Connect = db
		log.Println("Database connected...")
	}

	if err := Connect.Ping(); err != nil {
		log.Println("Database failed to ping")
		return err
	}

	return nil
}
