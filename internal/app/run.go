package app

import (
	"github.com/sirupsen/logrus"
	"www.github.com/kennnyz/avitochallenge/internal/config"
	httpdelivery "www.github.com/kennnyz/avitochallenge/internal/delivery/http"
	"www.github.com/kennnyz/avitochallenge/internal/repository"
	"www.github.com/kennnyz/avitochallenge/internal/server"
	service2 "www.github.com/kennnyz/avitochallenge/internal/service"
)

func Run() {
	cfg, err := config.ReadConfig()
	if err != nil {
		logrus.Panic("couldn't read config")
	}

	repos, err := repository.NewRepository(cfg)
	userSegmentService := service2.NewUserSegment(repos)
	handler := httpdelivery.NewHandler(userSegmentService)
	httpServer := server.NewHTTPServer(cfg.ServerAddr, handler.Init())

	logrus.Println("Server is listening..." + cfg.ServerAddr)
	err = httpServer.Run()
	if err != nil {
		logrus.Panic("Couldn't run server")
	}
}
