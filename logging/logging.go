package logging

import "go.uber.org/zap"

var Logger *zap.Logger

func init() {
	Logger = zap.Must(zap.NewProduction())
}
