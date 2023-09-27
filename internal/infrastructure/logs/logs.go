package logs

import (
	"github.com/KuYaki/waffler_server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	Debug  = "debug"
	Info   = "info"
	Empty  = ""
	Warn   = "warn"
	Error  = "error"
	Dpanic = "dpanic"
	Panic  = "panic"
	Fatal  = "fatal"
)

func NewLogger(conf config.AppConf) (*zap.Logger, error) {
	levels := map[string]zapcore.Level{
		Debug:  zapcore.DebugLevel,
		Info:   zapcore.InfoLevel,
		Empty:  zapcore.InfoLevel,
		Warn:   zapcore.WarnLevel,
		Error:  zapcore.ErrorLevel,
		Dpanic: zapcore.DPanicLevel,
		Panic:  zapcore.PanicLevel,
		Fatal:  zapcore.FatalLevel,
	}

	zapConf := zap.NewProductionConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConf.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	atom := zap.NewAtomicLevelAt(levels[conf.Logger.Level])
	zapConf.Level = atom

	logFile, err := createLogFileIfNotExists(conf.Logger.LogPath)
	if err != nil {
		return nil, err
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConf.EncoderConfig),
		zapcore.AddSync(logFile),
		atom,
	)

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConf.EncoderConfig),
		os.Stdout,
		atom,
	)

	logger := zap.New(zapcore.NewTee(fileCore, consoleCore), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger.Named(conf.AppName), nil
}
