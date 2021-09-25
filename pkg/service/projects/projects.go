package projects_service

import (
	"context"

	"github.com/droplez/droplez-go-proto/pkg/common"
	proto_projects "github.com/droplez/droplez-go-proto/pkg/studio/projects"
	projects_repo "github.com/droplez/droplez-studio/pkg/repo/projects"
	"github.com/droplez/droplez-studio/third_party/postgres"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var repo projects_repo.ProjectStore

var initRepo = func(ctx context.Context) projects_repo.ProjectStore {
	if repo == nil {
		repo = projects_repo.ProjectRepo{
			Pool: postgres.Pool(ctx),
		}
	}
	return repo
}

// Create a new project
func Create(ctx context.Context, in *proto_projects.ProjectMeta) (out *proto_projects.ProjectInfo, err error) {
	// Prepare repo layer
	repo := initRepo(ctx)

	// Create new project
	out = &proto_projects.ProjectInfo{
		Metadata: in,
		Id:       uuid.New().String(),
	}

	code, err := repo.CreateProject(ctx, out)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return
}

// Update a project
func Update(ctx context.Context, in *proto_projects.ProjectInfo) (*proto_projects.ProjectInfo, error) {
	// Prepare repo layer
	repo := initRepo(ctx)

	// Update project
	code, err := repo.UpdateProject(ctx, in)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return in, nil
}

func Get(ctx context.Context, in *proto_projects.ProjectId) (*proto_projects.ProjectInfo, error) {
	// Prepare repo layer
	repo := initRepo(ctx)

	// Update project
	project, code, err := repo.GetProject(ctx, in)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return project, nil
}

// Update application
func Delete(ctx context.Context, in *proto_projects.ProjectInfo) (*common.EmptyMessage, error) {
	repo := initRepo(ctx)
	projectID := &proto_projects.ProjectId{
		Id: in.GetId(),
	}
	projectGotten, code, err := repo.GetProject(ctx, projectID)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	if projectGotten.Metadata.Name == in.Metadata.Name {
		code, err := repo.DeleteProject(ctx, projectID)
		if err != nil {
			return nil, status.Error(code, err.Error())
		}
		return &common.EmptyMessage{}, nil
	}

	return nil, status.Error(codes.NotFound, "Project with this ID has another Name")
}

func List(ctx context.Context, stream proto_projects.Projects_ListServer, options *proto_projects.ListOptions) error {
	repo := initRepo(ctx)
	code, err := repo.ListUsers(ctx, stream, options)
	if err != nil {
		return status.Error(code, err.Error())
	}
	return nil
}