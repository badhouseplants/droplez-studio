package api

import (
	"context"

	"github.com/droplez/droplez-go-proto/pkg/common"
	"github.com/droplez/droplez-go-proto/pkg/studio/projects"
	"github.com/droplez/droplez-studio/pkg/service"
	"github.com/droplez/droplez-studio/tools/logger"
	"google.golang.org/grpc"
)

type projectsGrpcImpl struct {
	projects.UnimplementedProjectsServer
}

func RegisterProjectsServer(grpcServer *grpc.Server) {
	projects.RegisterProjectsServer(grpcServer, &projectsGrpcImpl{})
}

func (s projectsGrpcImpl) Create(ctx context.Context, in *projects.ProjectMeta) (*projects.ProjectInfo, error) {
	logger.EndpointHit(ctx)
	return service.ProjectCreate(ctx, in)
}

func (s projectsGrpcImpl) Update(ctx context.Context, in *projects.ProjectInfo) (*projects.ProjectInfo, error) {
	logger.EndpointHit(ctx)
	return service.ProjectUpdate(ctx, in)
}

func (s projectsGrpcImpl) Get(ctx context.Context, in *projects.ProjectId) (*projects.ProjectInfo, error) {
	logger.EndpointHit(ctx)
	return service.ProjectGet(ctx, in)
}

func (s projectsGrpcImpl) Delete(ctx context.Context, in *projects.ProjectInfo) (*common.EmptyMessage, error) {
	logger.EndpointHit(ctx)
	return service.ProjectDelete(ctx, in)
}

func (s projectsGrpcImpl) List(in *projects.ListOptions, stream projects.Projects_ListServer) (err error) {
	logger.EndpointHit(stream.Context())
	err = service.ProjectsList(stream.Context(), stream, in)
	if err != nil {
		return
	}
	return
}
