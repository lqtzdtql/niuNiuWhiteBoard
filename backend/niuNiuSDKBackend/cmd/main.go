package main

import (
	"context"
	"net/http"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/config"
	"niuNiuSDKBackend/internal/server"
	"niuNiuSDKBackend/internal/service"
	"time"
)

func main() {
	log.InitLogger(config.GetConfig().Log.Path, config.GetConfig().Log.Level)
	log.Logger.Info("config", log.Any("config", config.GetConfig()))

	log.Logger.Info("start server", log.String("start", "start web sever..."))

	newRouter := service.NewRouter()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.MyServer.Start(ctx)

	s := &http.Server{
		Addr:           ":8888",
		Handler:        newRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if nil != err {
		log.Logger.Error("server error", log.Any("serverError", err))
	}
}
