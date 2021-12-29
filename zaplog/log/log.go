package log

import (
	"go.uber.org/zap"
	"time"
)

var ZapLogger *zap.Logger

func init() {
	ZapLogger, _ := zap.NewProduction()
	sugar := ZapLogger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "www.baidu.com",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", "www.baidu.com")

}
