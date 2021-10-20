package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/droplez/droplez-go-proto/pkg/studio/projects"
	"github.com/droplez/droplez-studio/tools/logger"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
)

type ProjectRepo struct {
	Pool *pgxpool.Conn
}

func (r ProjectRepo) CreateProject(ctx context.Context, project *projects.ProjectInfo) (code codes.Code, err error) {
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

func (r ProjectRepo) UpdateProject(ctx context.Context, project *projects.ProjectInfo) (code codes.Code, err error) {
	const sql = "UPDATE projects SET name=$2, description=$3, public=$4, bpm=$5, key=$6, genre=$7, daw=$8 WHERE id=$1 RETURNING *"

	var log = logger.GetGrpcLogger(ctx)

	_, err = r.Pool.Exec(ctx, sql,
		project.Id, project.Metadata.Name,
		project.Metadata.Description, project.Metadata.Public,
		project.Metadata.Bpm, project.Metadata.Key,
		project.Metadata.Daw.String(),
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

func (r ProjectRepo) GetProject(ctx context.Context, projectID *projects.ProjectId) (*projects.ProjectInfo, codes.Code, error) {
	const sql = "SELECT name, description, public, bpm, key, genre, daw FROM projects WHERE id = $1"

	var log = logger.GetGrpcLogger(ctx)
	var projectMeta = &projects.ProjectMeta{}
	var daw string

	err := r.Pool.QueryRow(ctx, sql, projectID.GetId()).Scan(
		&projectMeta.Name, &projectMeta.Description, &projectMeta.Public,
		&projectMeta.Bpm, &projectMeta.Key, &projectMeta.Genre,
		&daw,
	)
	projectMeta.Daw = projects.DAW(projects.DAW_value[daw])

	project := &projects.ProjectInfo{
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

func (r ProjectRepo) DeleteProject(ctx context.Context, projectID *projects.ProjectId) (codes.Code, error) {
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

func (r ProjectRepo) ListProjects(ctx context.Context, stream projects.Projects_ListServer, opt *projects.ListOptions) (codes.Code, error) {
	const sql = "SELECT id, name, description, public, bpm, key, genre, daw FROM projects LIMIT $1 OFFSET $2"
	var (
		log         = logger.GetGrpcLogger(ctx)
		project     = &projects.ProjectInfo{}
		projectMeta = &projects.ProjectMeta{}
		projectID   = &projects.ProjectId{}
		daw         string
	)

	rows, err := r.Pool.Query(ctx, sql, opt.GetPaging().GetCount(), opt.GetPaging().GetPage())
	if err != nil {
		log.Error(err)
		return codes.Internal, err
	}

	for rows.Next() {
		err = rows.Scan(
			&projectID.Id, &projectMeta.Name,
			&projectMeta.Description, &projectMeta.Public,
			&projectMeta.Bpm, &projectMeta.Key,
			&projectMeta.Genre, &daw,
		)
		projectMeta.Daw = projects.DAW(projects.DAW_value[daw])
		if err != nil {
			log.Error(err)
			return codes.Internal, err
		}
		project.Id = projectID
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
