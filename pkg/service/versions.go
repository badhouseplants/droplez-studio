package service

import (
	"context"

	"github.com/droplez/droplez-go-proto/pkg/common"
	"github.com/droplez/droplez-go-proto/pkg/studio/versions"
	"github.com/droplez/droplez-studio/pkg/repo"
	"github.com/droplez/droplez-studio/third_party/postgres"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type VersionStore interface {
	CreateVersion(ctx context.Context, in *versions.VersionInfo) (codes.Code, error)
	UpdateVersion(ctx context.Context, in *versions.VersionInfo) (codes.Code, error)
	GetVersions(ctx context.Context, in *versions.VersionId) (*versions.VersionInfo, codes.Code, error)
	ListVersions(ctx context.Context, stream versions.Versions_ListServer, options *versions.ListOptions) (codes.Code, error)
}

var versionStore VersionStore

var initVersionsRepo = func(ctx context.Context) VersionStore {
	if versionStore == nil {
		versionStore = repo.VersionRepo{
			Pool: postgres.Pool(ctx),
		}
	}
	return versionStore
}

func VersionCreate(ctx context.Context, in *versions.VersionMeta) (*versions.VersionInfo, error) {
	repo := initVersionsRepo(ctx)

	out := &versions.VersionInfo{
		Id: &versions.VersionId{
			Id: uuid.New().String(),
		},
		Metadata: in,
	}
	out.Metadata.UploadedAt = timestamppb.Now()
	code, err := repo.CreateVersion(ctx, out)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return out, nil
}

func VersionUpdate(ctx context.Context, in *versions.VersionInfo) (*versions.VersionInfo, error) {
	repo := initVersionsRepo(ctx)
	out := in
	out.Metadata.UploadedAt = timestamppb.Now()
	code, err := repo.UpdateVersion(ctx, out)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}
	return out, nil
}

func VersionDelete(ctx context.Context, in *versions.VersionInfo) (*common.EmptyMessage, error) {
	return nil, status.Error(codes.Unimplemented, "delete versions is not allowed yet")
}

func VersionGet(ctx context.Context, in *versions.VersionId) (*versions.VersionInfo, error) {
	repo := initVersionsRepo(ctx)

	out, code, err := repo.GetVersions(ctx, in)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}
	return out, nil
}

func VersionsList(ctx context.Context, stream versions.Versions_ListServer, options *versions.ListOptions) error {
	repo := initVersionsRepo(ctx)
	code, err := repo.ListVersions(ctx, stream, options)
	if err != nil {
		return status.Error(code, err.Error())
	}
	return nil

}
