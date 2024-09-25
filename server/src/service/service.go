package service

import (
	"quizserver/config"
	"quizserver/logger"
)

// Service ...
type Service struct {
	Logger *logger.Log
	Config *config.Config
}

func NewService(
	Logger *logger.Log,
	Config *config.Config,
) Service {
	return Service{
		Logger: Logger,
		Config: Config,
	}
}

func (s *Service) Close() {

}
