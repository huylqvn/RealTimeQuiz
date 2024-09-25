package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"quizserver/config"
	"quizserver/logger"
	"quizserver/src/endpoints"
	serviceHttp "quizserver/src/http"
	"quizserver/src/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	log := logger.Get().
		SetCaller().
		SetTimestamp().
		SetLevel(cfg.LogLevel).
		SetSource("quizserver")
	if cfg.LogFile == "true" {
		log.SetFile()
	}
	log.Info("LoadConfig", "ok")

	if cfg.Port == "" {
		cfg.Port = "7000"
	}

	port := fmt.Sprintf(":%s", cfg.Port)
	var (
		httpAddr = flag.String("addr", port, "HTTP listen address")
	)

	var s service.Service
	var h http.Handler

	log.Info("database", "connected to db")

	log.Info("redis", "connected to redis")

	s = service.NewService(log, cfg)
	defer s.Close()

	s.Logger.Info("deployVersion", s.Config.Version)

	h = serviceHttp.NewHTTPHandler(
		s,
		endpoints.MakeServerEndpoints(&s),
		log.GetLogger(),
		true,
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		log.Info("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	log.Log("exit", <-errs)
}
