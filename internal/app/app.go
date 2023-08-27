package app

import (
	"github.com/sirupsen/logrus"
	"www.github.com/kennnyz/avitochallenge/internal/config"
	httpdelivery "www.github.com/kennnyz/avitochallenge/internal/delivery/http"
	"www.github.com/kennnyz/avitochallenge/internal/repository"
	"www.github.com/kennnyz/avitochallenge/internal/server"
	service2 "www.github.com/kennnyz/avitochallenge/internal/service"
	"www.github.com/kennnyz/avitochallenge/pkg/database"
)

func Run() {
	cfg, err := config.ReadConfig()
	if err != nil {
		logrus.Panic("couldn't read config")
	}

	db, err := database.NewClient(cfg.DB.Dsn)
	if err != nil {
		logrus.Panic(err)
	}
	repos := repository.NewUserSegmentRepository(db)
	userSegmentService := service2.NewUserSegment(repos)
	handler := httpdelivery.NewHandler(userSegmentService)
	httpServer := server.NewHTTPServer(cfg.ServerAddr, handler.Init())

	logrus.Println("Server is listening..." + cfg.ServerAddr)
	err = httpServer.Run()
	if err != nil {
		logrus.Panic("Couldn't run server")
	}
}
