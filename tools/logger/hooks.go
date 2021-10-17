package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

func EndpointHit(ctx context.Context) *logrus.Entry{
	log := GetGrpcLogger(ctx)
	log.Info("endpoint hit")
	return logger
}

