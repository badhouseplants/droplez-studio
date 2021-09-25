package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

func EnpointHit(ctx context.Context) *logrus.Entry{
	log := GetGrpcLogger(ctx)
	log.Info("enpoint hit")
	return logger
}

