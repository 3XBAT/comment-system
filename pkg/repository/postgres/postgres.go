package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	postsTable = "posts"
	commentsTable = "comments"
)

type Config struct {
	Port     string
	Username string
	Host     string
	DBName   string
	Password string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error){
	db, err := sqlx.Open("postgres", fmt.Sprintf("port=%s user=%s host=%s dbname=%s password=%s sslmode=%s",
	cfg.Port, cfg.Username, cfg.Host, cfg.DBName, cfg.Password, cfg.SSLMode))
	logrus.Printf("port=%s user=%s host=%s dbname=%s password=%s sslmode=%s",
	cfg.Port, cfg.Username, cfg.Host, cfg.DBName, cfg.Password, cfg.SSLMode)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	
	logrus.Println("db is successful created")

	return db, nil
} 
