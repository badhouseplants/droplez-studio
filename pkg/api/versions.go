package api

import (
	"context"

	"github.com/droplez/droplez-go-proto/pkg/common"
	"github.com/droplez/droplez-go-proto/pkg/studio/versions"
	"github.com/droplez/droplez-studio/pkg/service"
	"github.com/droplez/droplez-studio/tools/logger"
	"google.golang.org/grpc"
)

type versionsGrpcImpl struct {
	versions.UnimplementedVersionsServer
}

func RegisterVersionsServer(grpcServer *grpc.Server) {
	versions.RegisterVersionsServer(grpcServer, &versionsGrpcImpl{})
}

func (s versionsGrpcImpl) Create(ctx context.Context, in *versions.VersionMeta) (*versions.VersionInfo, error) {
	logger.EndpointHit(ctx)
	return service.VersionCreate(ctx, in)
}

func (s versionsGrpcImpl) Update(ctx context.Context, in *versions.VersionInfo) (*versions.VersionInfo, error) {
	logger.EndpointHit(ctx)
	return service.VersionUpdate(ctx, in)
}

func (s versionsGrpcImpl) Get(ctx context.Context, in *versions.VersionId) (*versions.VersionInfo, error) {
	logger.EndpointHit(ctx)
	return service.VersionGet(ctx, in)
}

func (s versionsGrpcImpl) List(in *versions.ListOptions, stream versions.Versions_ListServer) (err error) {
	logger.EndpointHit(stream.Context())
	err = service.VersionsList(stream.Context(), stream, in)
	if err != nil {
		return
	}
	return
}

func (s versionsGrpcImpl) Delete(ctx context.Context, in *versions.VersionInfo) (*common.EmptyMessage, error) {
	logger.EndpointHit(ctx)
	return service.VersionDelete(ctx, in)

}
