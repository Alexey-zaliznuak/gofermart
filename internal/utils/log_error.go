package utils

import (
	"github.com/Alexey-zaliznuak/gofermart/internal/logger"
	"go.uber.org/zap"
)

func LogErrorWrapper(err error) {
	if err != nil {
		logger.Log.Error("Error wrapped", zap.Error(err))
	}
}
