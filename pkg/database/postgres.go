package database

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	// pgx import
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewClient(dsn string) (*sql.DB, error) {
	counts := 0

	for {
		connection, err := openDB(dsn)
		if err != nil {
			logrus.Println("postgres not yet ready...", err)
		} else {
			logrus.Println("connected to database!")
			return connection, nil
		}

		if counts > 10 {
			return nil, fmt.Errorf("DB its sleep. ")
		}

		logrus.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
