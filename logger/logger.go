package logger

import (
	"github.com/pkg/errors"
	"myGin/global"
	"myGin/settings"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化Logger
func Init() error {
	logFilePath := path.Join(settings.AppCfg.RootDir, "logs", settings.LogCfg.Filename)
	writeSyncer := getLogWriter(
		logFilePath,
		settings.LogCfg.MaxSize,
		settings.LogCfg.MaxBackup,
		settings.LogCfg.MaxAge)

	encoder := getEncoder()
	var l = new(zapcore.Level)
	if err := l.UnmarshalText([]byte(settings.LogCfg.Level)); err != nil {
		return errors.Wrapf(err, "unmarshal txt level failed")
	}

	var core zapcore.Core
	// 开发者模式
	if settings.AppCfg.RunMode == global.DebugModel || settings.AppCfg.RunMode == global.InnerTestModel {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	//lg := zap.New(core, zap.AddCaller())
	lg := zap.New(core)
	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}