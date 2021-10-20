package repo

import (
	"context"
	"errors"
	"time"

	"github.com/droplez/droplez-go-proto/pkg/studio/versions"
	"github.com/droplez/droplez-studio/tools/logger"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type VersionRepo struct {
	Pool *pgxpool.Conn
}

func (r VersionRepo) CreateVersion(ctx context.Context, version *versions.VersionInfo) (code codes.Code, err error) {
	const sql = "INSERT INTO versions (id, version, project_id, object_name, message, uploaded_at) VALUES ($1, $2, $3, $4, $5, $6)"
	log := logger.GetGrpcLogger(ctx)

	_, err = r.Pool.Exec(ctx, sql,
		version.GetId().GetId(), version.GetMetadata().GetVersion(),
		version.GetMetadata().GetProjectId(), version.GetMetadata().GetObjectName(),
		version.GetMetadata().Message, version.GetMetadata().GetUploadedAt().AsTime(),
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

func (r VersionRepo) UpdateVersion(ctx context.Context, version *versions.VersionInfo) (codes.Code, error) {
	const sql = "UPDATE versions SET version=$2, project_id=$3, object_name=$4, message=$5, uploaded_at=$6 WHERE id=$1 RETURNING *"
	log := logger.GetGrpcLogger(ctx)

	_, err := r.Pool.Exec(ctx, sql,
		version.GetId().GetId(), version.GetMetadata().GetVersion(),
		version.GetMetadata().GetProjectId(), version.GetMetadata().GetObjectName(),
		version.GetMetadata().Message, version.GetMetadata().GetUploadedAt().AsTime(),
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

func (r VersionRepo) GetVersions(ctx context.Context, in *versions.VersionId) (*versions.VersionInfo, codes.Code, error) {
	const sql = "SELECT id, version, project_id, object_name, message, uploaded_at FROM versions WHERE id=$1"
	var timestamp  time.Time
	var log = logger.GetGrpcLogger(ctx)
	version := &versions.VersionInfo{
		Id: &versions.VersionId{},
		Metadata: &versions.VersionMeta{},
	}


	err := r.Pool.QueryRow(ctx, sql, in.GetId()).Scan(
		&version.Id.Id, &version.Metadata.Version,
		&version.Metadata.ProjectId, &version.Metadata.ObjectName,
		&version.Metadata.Message, &timestamp,
	)

	version.Metadata.UploadedAt = &timestamppb.Timestamp{
		Seconds: timestamp.Unix(),
		Nanos: timestamppb.Now().Nanos,
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, codes.NotFound, errProjectNotFoundByID(version.GetId().GetId())
		} else {
			log.Error(err)
			return nil, codes.Internal, err
		}
	}

	return version, codes.OK, nil
}

func (r VersionRepo) ListVersions(ctx context.Context, stream versions.Versions_ListServer, opt *versions.ListOptions) (codes.Code, error) {
	const sql = "SELECT id, version, project_id, object_name, message, uploaded_at FROM versions LIMIT $1 OFFSET $2"
	var timestamp  time.Time
	var log = logger.GetGrpcLogger(ctx)
	version := &versions.VersionInfo{
		Id: &versions.VersionId{},
		Metadata: &versions.VersionMeta{},
	}
	
	rows, err := r.Pool.Query(ctx, sql, opt.GetPaging().GetCount(), opt.GetPaging().GetPage())
	if err != nil {
		log.Error(err)
		return codes.Internal, err
	}

	for rows.Next() {
		err = rows.Scan(
			&version.Id.Id, &version.Metadata.Version,
			&version.Metadata.ProjectId, &version.Metadata.ObjectName,
			&version.Metadata.Message, &timestamp,	
		)
		version.Metadata.UploadedAt = &timestamppb.Timestamp{
			Seconds: timestamp.Unix(),
			Nanos: timestamppb.Now().Nanos,
		}
	
	
		if err := stream.Send(version); err != nil {
			log.Error(err)
			return codes.Internal, err
		}
	}
	return codes.OK, nil

}