package service

import "go.uber.org/zap"

func NewLogger() *zap.SugaredLogger {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	return sugar
}
