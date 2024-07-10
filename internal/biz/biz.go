package biz

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/google/wire"
	"google.golang.org/grpc/metadata"
	v1 "helloworld/api/helloworld/v1"
	"helloworld/internal/data/ent"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase)

var (
	LoginNotExistsError = errors.New(401, v1.ErrorReason_NOT_AUTHORIZED.String(), "")
	ErrRouteNotActual   = errors.New(410, v1.ErrorReason_ROUTE_NOT_ACTUAL.String(), "")
	ErrRouteInProcess   = errors.New(202, v1.ErrorReason_ROUTE_IN_PROCESS.String(), "")
)

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	GetPasswordHash(ctx context.Context, login string) (uint64, string, error)
	CreateSession(ctx context.Context, userId uint64) (string, error)
	GetUserId(ctx context.Context, token string) (uint64, error)
	DeleteAllSession(ctx context.Context, userId uint64) error
	InsertAsset(ctx context.Context, userId uint64, assetName string, file []byte) error
	GetAsset(ctx context.Context, assetName string) ([]byte, error)
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *GreeterUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	userId, passwordHash, err := uc.repo.GetPasswordHash(ctx, req.Login)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, LoginNotExistsError
		}
		return nil, err
	}
	hash := md5.New()
	hash.Write([]byte(req.Password))
	passwordHashIn := hex.EncodeToString(hash.Sum(nil))
	if passwordHashIn != passwordHash {
		return nil, LoginNotExistsError
	}
	token, err := uc.repo.CreateSession(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func getToken(ctx context.Context) (string, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", LoginNotExistsError
	}
	auth := meta["authorization"][0]
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", LoginNotExistsError
	}
	return auth[7:], nil
}

func (uc *GreeterUsecase) UploadAsset(ctx context.Context, req *v1.AssetRequest) (*v1.StatusReply, error) {
	token, err := getToken(ctx)
	if err != nil {
		return nil, err
	}
	userId, err := uc.repo.GetUserId(ctx, token)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.InsertAsset(ctx, userId, req.AssetName, req.Data); err != nil {
		return nil, err
	}
	return &v1.StatusReply{
		Status: "ok",
	}, nil
}

func (uc *GreeterUsecase) Get(ctx context.Context, req *v1.AssetRequest) (*v1.GetReply, error) {
	token, err := getToken(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := uc.repo.GetUserId(ctx, token); err != nil {
		return nil, err
	}
	data, err := uc.repo.GetAsset(ctx, req.AssetName)
	if err != nil {
		return nil, err
	}
	return &v1.GetReply{
		Data: data,
	}, nil
}
