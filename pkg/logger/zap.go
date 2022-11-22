package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"top-ping/pkg/utils"
)

const (
	devProfile  = "dev"
	defaultFile = "default.log"
	errorFile   = "error.log"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func NewZapLogger(profile string, config *Config) *zap.Logger {
	if ok, _ := utils.PathExists(config.Dir); !ok {
		_ = os.Mkdir(config.Dir, os.ModePerm)
	}

	var cores []zapcore.Core
	cores = append(cores, newDefaultCore(profile, config))
	cores = append(cores, newErrorCore(profile, config))

	combinedCore := zapcore.NewTee(cores...)

	var options []zap.Option
	options = append(options, zap.AddCaller())
	options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	options = append(options, zap.AddCallerSkip(1))

	return zap.New(combinedCore, options...)
}

func newDefaultCore(profile string, config *Config) zapcore.Core {
	defaultWriter := newLogWriter(config.Dir+"/"+defaultFile, config)
	configLogLevel := getLoggerLevel(config.Level)
	defaultLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return configLogLevel <= lvl && lvl <= zapcore.FatalLevel
	})

	encoderCfg := newEncoderConfig(profile)
	var core zapcore.Core
	var writeSyncer zapcore.WriteSyncer
	if profile == devProfile {
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(defaultWriter))
	} else {
		writeSyncer = zapcore.AddSync(defaultWriter)
	}

	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writeSyncer,
		defaultLevel,
	)

	return core
}

func newErrorCore(profile string, config *Config) zapcore.Core {
	errorWriter := newLogWriter(config.Dir+"/"+errorFile, config)
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	encoderCfg := newEncoderConfig(profile)
	var core zapcore.Core
	var writeSyncer zapcore.WriteSyncer
	if profile == devProfile {
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(errorWriter))
	} else {
		writeSyncer = zapcore.AddSync(errorWriter)
	}

	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writeSyncer,
		errorLevel,
	)

	return core
}

func newLogWriter(filename string, config *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func newEncoderConfig(profile string) zapcore.EncoderConfig {
	var encoderCfg zapcore.EncoderConfig
	if profile == devProfile {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return encoderCfg
}

func getLoggerLevel(logLevel string) zapcore.Level {
	level, exist := loggerLevelMap[logLevel]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}
