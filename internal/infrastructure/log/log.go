package log

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.SugaredLogger
var logFile *os.File

const FileRWPerm = os.FileMode(0o666)

func init() {
	if isRunAsGenerate() || testing.Testing() {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic("cannot initialize zap logger")
		}

		Zap = logger.Sugar()
		return
	}

	// Создаем файл для записи логов
	file, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, FileRWPerm)
	logFile = file
	if err != nil {
		fmt.Printf("Failed to open log file: %s", err)
	}
	// Создаем два writer: один для stdout, другой для файла.
	consoleWriter := zapcore.AddSync(os.Stdout)
	fileWriter := zapcore.AddSync(logFile)
	// Выбираем формат вывода.
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	// Уровень логирования (например, DebugLevel).
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleWriter, zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriter, zapcore.DebugLevel),
	)
	// Создаем логгер.
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	// Используем SugaredLogger для удобства.
	Zap = logger.Sugar()
}

func Close() error {
	Zap.Info("Logger closing")
	err := Zap.Sync()
	if err != nil {
		fmt.Printf("failed zap logger sync: %s", err.Error())
	}

	err = logFile.Close()
	if err != nil {
		return fmt.Errorf("failed close log file: %w", err)
	}
	return nil
}

func isRunAsGenerate() bool {
	return os.Getenv("GOFILE") != ""
}
