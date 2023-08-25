package repository

import (
	"github.com/sirupsen/logrus"
	"www.github.com/kennnyz/avitochallenge/internal/config"
	"www.github.com/kennnyz/avitochallenge/internal/repository/postgres"
	"www.github.com/kennnyz/avitochallenge/internal/service"
	"www.github.com/kennnyz/avitochallenge/pkg/database"
)

func NewRepository(cfg *config.Config) (service.UserSegmentRepository, error) {
	db, err := database.NewClient(cfg.DB.Dsn)
	if err != nil {
		logrus.Panic(err)
	}
	return postgres.NewRepository(db), nil
}
