package logger

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

func GetSugar() *zap.SugaredLogger {
	return sugar
}

func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()
	sugar = logger.Sugar()

	return nil
}
