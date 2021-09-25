package projects_repo

import (
	"context"
	"errors"
	"fmt"

	proto_projects "github.com/droplez/droplez-go-proto/pkg/studio/projects"
	"github.com/droplez/droplez-studio/tools/logger"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
)

type ProjectStore interface {
	CreateProject(context.Context, *proto_projects.ProjectInfo) (codes.Code, error)
	UpdateProject(context.Context, *proto_projects.ProjectInfo) (codes.Code, error)
	GetProject(context.Context, *proto_projects.ProjectId) (*proto_projects.ProjectInfo, codes.Code, error)
	DeleteProject(context.Context, *proto_projects.ProjectId) (codes.Code, error)
	ListUsers(context.Context, proto_projects.Projects_ListServer, *proto_projects.ListOptions) (codes.Code, error)
}

type ProjectRepo struct {
	Pool *pgxpool.Conn
}

func (r ProjectRepo) CreateProject(ctx context.Context, project *proto_projects.ProjectInfo) (code codes.Code, err error) {
	const sql = `INSERT INTO projects 
								(id, name, daw, description, public, bpm, key, genre) 
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	var log = logger.GetGrpcLogger(ctx)
	_, err = r.Pool.Exec(ctx, sql,
		project.Id.Id, project.Metadata.Name,
		project.Metadata.Daw.String(), project.Metadata.Description,
		project.Metadata.Public, project.Metadata.Bpm,
		project.Metadata.Key, project.Metadata.Genre,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return codes.AlreadyExists, err
			default:
				log.Error(err)
				return codes.Internal, err
			}
		}
		log.Error(err)
		return codes.Internal, err
	}

	return codes.OK, nil
}

func (r ProjectRepo) UpdateProject(ctx context.Context, project *proto_projects.ProjectInfo) (code codes.Code, err error) {
	const sql = "UPDATE projects SET name=$2, description=$3, public=$4, bpm=$5, key=$6, genre=$7, mood=$8 WHERE id=$1 RETURNING *"

	var log = logger.GetGrpcLogger(ctx)

	_, err = r.Pool.Exec(ctx, sql,
		project.Id, project.Metadata.Name,
		project.Metadata.Description, project.Metadata.Public,
		project.Metadata.Bpm, project.Metadata.Key,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return codes.AlreadyExists, err
			default:
				log.Error(err)
				return codes.Internal, err
			}
		}
	}

	return codes.OK, nil
}

func (r ProjectRepo) GetProject(ctx context.Context, projectID *proto_projects.ProjectId) (*proto_projects.ProjectInfo, codes.Code, error) {
	const sql = "SELECT name, description, public, bpm, key, genre, mood,  FROM projects WHERE id = $1"

	var log = logger.GetGrpcLogger(ctx)
	var projectMeta = &proto_projects.ProjectMeta{}

	err := r.Pool.QueryRow(ctx, sql, projectID.GetId()).Scan(
		&projectMeta.Name, &projectMeta.Description,
		&projectMeta.Public, &projectMeta.Bpm, &projectMeta.Key, &projectMeta.Genre,
	)

	project := &proto_projects.ProjectInfo{
		Metadata: projectMeta,
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, codes.NotFound, errProjectNotFoundByID(projectID.GetId())
		} else {
			log.Error(err)
			return nil, codes.Internal, err
		}
	}

	return project, codes.OK, nil
}

func (r ProjectRepo) DeleteProject(ctx context.Context, projectID *proto_projects.ProjectId) (codes.Code, error) {
	const sql = "DELETE FROM projects WHERE id = $1"

	log := logger.GetGrpcLogger(ctx)

	tag, err := r.Pool.Exec(ctx, sql, projectID.GetId())
	if tag.RowsAffected() == 0 {
		return codes.NotFound, errProjectNotFoundByID(projectID.GetId())
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return codes.NotFound, errProjectNotFoundByID(projectID.GetId())
		} else {
			log.Error(err)
			return codes.Internal, err
		}
	}

	return codes.OK, nil
}

func (r ProjectRepo) ListUsers(ctx context.Context, stream proto_projects.Projects_ListServer, opt *proto_projects.ListOptions) (codes.Code, error) {
	const sql = "SELECT id, name, description, public, bpm, key, genre, mood FROM projects LIMIT $1 OFFSET $2"
	var (
		log         = logger.GetGrpcLogger(ctx)
		project     = &proto_projects.ProjectInfo{}
		projectMeta = &proto_projects.ProjectMeta{}
	)

	rows, err := r.Pool.Query(ctx, sql, opt.GetPaging().GetCount(), opt.GetPaging().GetPage())
	if err != nil {
		log.Error(err)
		return codes.Internal, err
	}

	for rows.Next() {
		err = rows.Scan(
			&project.Id, &projectMeta.Name,
			&projectMeta.Description, &projectMeta.Public,
			&projectMeta.Bpm, &projectMeta.Key,
		)
		if err != nil {
			log.Error(err)
			return codes.Internal, err
		}
		project.Metadata = projectMeta
		if err := stream.Send(project); err != nil {
			log.Error(err)
			return codes.Internal, err
		}
	}
	return codes.OK, nil

}

//Local errors
var (
	errProjectNotFoundByID = func(id string) error {
		return fmt.Errorf("project with this id can not be found: %s", id)
	}
)
