package main

import (
	"fmt"

	"go-spider/zaplog"
)

const logPath = "./logs"

var log = zaplog.Logger

func initLog() {
	_logPath := fmt.Sprintf("%s/%s", logPath, "info.log")
	errLogPath := fmt.Sprintf("%s/%s", logPath, "info-err.log")
	_logger := zaplog.NewLogger(
		_logPath,
		zaplog.WithErrFilePath(errLogPath),
		zaplog.WithFileMaxAgeDay(15),
	)
	log = _logger
}
