package service

import (
	"context"

	"github.com/droplez/droplez-go-proto/pkg/common"
	"github.com/droplez/droplez-go-proto/pkg/studio/projects"
	"github.com/droplez/droplez-studio/pkg/repo"
	"github.com/droplez/droplez-studio/third_party/postgres"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectStore interface {
	CreateProject(context.Context, *projects.ProjectInfo) (codes.Code, error)
	UpdateProject(context.Context, *projects.ProjectInfo) (codes.Code, error)
	GetProject(context.Context, *projects.ProjectId) (*projects.ProjectInfo, codes.Code, error)
	DeleteProject(context.Context, *projects.ProjectId) (codes.Code, error)
	ListProjects(context.Context, projects.Projects_ListServer, *projects.ListOptions) (codes.Code, error)
}

var projectStore ProjectStore

var initProjectRepo = func(ctx context.Context) ProjectStore {
	if projectStore == nil {
		projectStore = repo.ProjectRepo{
			Pool: postgres.Pool(ctx),
		}
	}
	return projectStore
}

// ProjectCreate a new project
func ProjectCreate(ctx context.Context, in *projects.ProjectMeta) (out *projects.ProjectInfo, err error) {
	// Prepare repo layer
	repo := initProjectRepo(ctx)

	// Create new project
	out = &projects.ProjectInfo{
		Metadata: in,
		Id: &projects.ProjectId{
			Id: uuid.New().String(),
		},
	}

	code, err := repo.CreateProject(ctx, out)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return
}

// ProjectUpdate a project
func ProjectUpdate(ctx context.Context, in *projects.ProjectInfo) (*projects.ProjectInfo, error) {
	// Prepare repo layer
	repo := initProjectRepo(ctx)

	// Update project
	code, err := repo.UpdateProject(ctx, in)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return in, nil
}

func ProjectGet(ctx context.Context, in *projects.ProjectId) (*projects.ProjectInfo, error) {
	// Prepare repo layer
	repo := initProjectRepo(ctx)

	// Update project
	project, code, err := repo.GetProject(ctx, in)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return project, nil
}

// Update application
func ProjectDelete(ctx context.Context, in *projects.ProjectInfo) (*common.EmptyMessage, error) {
	repo := initProjectRepo(ctx)
	projectID := &projects.ProjectId{
		Id: in.GetId().GetId(),
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

func ProjectsList(ctx context.Context, stream projects.Projects_ListServer, options *projects.ListOptions) error {
	repo := initProjectRepo(ctx)
	code, err := repo.ListProjects(ctx, stream, options)
	if err != nil {
		return status.Error(code, err.Error())
	}
	return nil
}
