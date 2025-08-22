package initialization

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"id-backend-grpc/internal/app/config"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
)

const (
	defaultEncoding   = "json"
	defaultLevel      = zapcore.InfoLevel
	defaultOutputPath = "id-backend.log"
)

func CreateLogger(cfg *config.LoggingConfig) (*zap.Logger, error) {
	return zap.NewProductionConfig().Build()
	//level := defaultLevel
	//output := defaultOutputPath
	//
	//if cfg != nil {
	//	if cfg.Level != "" {
	//		supportedLoggingLevels := map[string]zapcore.Level{
	//			debugLevel: zapcore.DebugLevel,
	//			infoLevel:  zapcore.InfoLevel,
	//			warnLevel:  zapcore.WarnLevel,
	//			errorLevel: zapcore.ErrorLevel,
	//		}
	//		var found bool
	//		if level, found = supportedLoggingLevels[cfg.Level]; !found {
	//			return nil, errors.New("logging level is incorrect")
	//		}
	//	}
	//
	//	if cfg.Output != "" {
	//		// TODO: need to create a
	//		// directory if it is missing
	//		output = cfg.Output
	//	}
	//}
	//
	//loggerCfg := zap.Config{
	//	Encoding:    defaultEncoding,
	//	Level:       zap.NewAtomicLevelAt(level),
	//	OutputPaths: []string{output},
	//}
	//
	//return loggerCfg.Build()
}
