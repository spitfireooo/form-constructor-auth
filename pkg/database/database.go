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
	UserTable = "users"
)

var Connect *sqlx.DB

func DatabaseInit(con DBConfig) error {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		con.Username, con.Password, con.Database, con.Host, con.Port, con.SSLMode,
	)

	if db, err := sqlx.Open("postgres", dsn); err != nil {
		return err
	} else {
		Connect = db
		log.Println("Database connect... Successfully")
		return nil
	}
}
