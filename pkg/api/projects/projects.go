package projects_api

import (
	"context"

	"github.com/droplez/droplez-go-proto/pkg/common"
	proto_projects "github.com/droplez/droplez-go-proto/pkg/studio/projects"
	service_projects "github.com/droplez/droplez-studio/pkg/service/projects"
	"github.com/droplez/droplez-studio/tools/logger"
	"google.golang.org/grpc"
)

type projectsGrpcImpl struct {
	proto_projects.UnimplementedProjectsServer
}

func Register(grpcServer *grpc.Server) {
	proto_projects.RegisterProjectsServer(grpcServer, &projectsGrpcImpl{})
}

func (s projectsGrpcImpl) Create(ctx context.Context, in *proto_projects.ProjectMeta) (*proto_projects.ProjectInfo, error) {
	logger.EndpointHit(ctx)
	return service_projects.Create(ctx, in)
}

func (s projectsGrpcImpl) Update(ctx context.Context, in *proto_projects.ProjectInfo) (*proto_projects.ProjectInfo, error) {
	logger.EndpointHit(ctx)
	return service_projects.Update(ctx, in)
}

func (s projectsGrpcImpl) Get(ctx context.Context, in *proto_projects.ProjectId) (*proto_projects.ProjectInfo, error) {
	logger.EndpointHit(ctx)
	return service_projects.Get(ctx, in)
}

func (s projectsGrpcImpl) Delete(ctx context.Context, in *proto_projects.ProjectInfo) (*common.EmptyMessage, error) {
	logger.EndpointHit(ctx)
	return service_projects.Delete(ctx, in)
}

func (s projectsGrpcImpl) List(in *proto_projects.ListOptions, stream proto_projects.Projects_ListServer) (err error) {
	logger.EndpointHit(stream.Context())
	err = service_projects.List(stream.Context(), stream, in)
	if err != nil {
		return
	}
	return
}
