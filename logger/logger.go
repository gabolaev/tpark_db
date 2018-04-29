package logger

import (
	"go.uber.org/zap"
)

var Instance, _ = zap.NewDevelopment()
