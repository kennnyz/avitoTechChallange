package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewPgxClient(dsn string) (*pgxpool.Pool, error) {
	counts := 0

	for {
		connection, err := openPgxDB(dsn)
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

func openPgxDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil
}
