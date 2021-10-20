package server

import (
	"fmt"
	"net"

	"github.com/droplez/droplez-studio/pkg/api"
	"github.com/droplez/droplez-studio/tools/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// grpc server options
var opts = []grpc.ServerOption{grpc.MaxRecvMsgSize(2147483648), setupGrpcUnaryOpts(), setupGrpcStreamOpts()}

// grpc server with services
var grpcServer = func() (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer(opts...)
	// Register services
	api.RegisterProjectsServer(grpcServer)
	api.RegisterVersionsServer(grpcServer)
	reflection.Register(grpcServer)
	return
}

// Serve starts grpc server
func Serve() (err error) {
	log := logger.GetServerLogger()

	// Preparing listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "9090")) //TODO: get port from env
	if err != nil {
		return err
	}

	// Start grpc server
	if err = grpcServer().Serve(listener); err != nil {
		log.Error(err)
		return err
	}

	// Shutdown
	return nil
}

func setupGrpcUnaryOpts() grpc.ServerOption {
	return grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logger.GrpcLogrusEntry, logger.GrpcLogrusOpts...),
		grpc_recovery.UnaryServerInterceptor(),
	)
}

func setupGrpcStreamOpts() grpc.ServerOption {
	return grpc_middleware.WithStreamServerChain(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.StreamServerInterceptor(logger.GrpcLogrusEntry, logger.GrpcLogrusOpts...),
		grpc_recovery.StreamServerInterceptor(),
	)
}
