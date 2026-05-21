package logger

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

// GetSugar returns the shared sugared logger instance.
//
// The logger should be initialized with InitLogger before use.
func GetSugar() *zap.SugaredLogger {
	return sugar
}

// InitLogger initializes the shared development logger instance.
func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()
	sugar = logger.Sugar()

	return nil
}
