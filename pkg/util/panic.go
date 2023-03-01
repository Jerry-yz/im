package util

import (
	"fmt"
	"learn-im/logger"
	"runtime"

	"go.uber.org/zap"
)

func Recover() {
	err := recover()
	if err != nil {
		logger.Logger.DPanic("recover panic", zap.Any("panic", err), zap.String("stack", GetStackInfo()))
	}
}

func GetStackInfo() string {
	buf := make([]byte, 4069)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}
