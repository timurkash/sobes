package service

import (
	"context"
	"github.com/google/wire"

	v1 "helloworld/api/helloworld/v1"
	"helloworld/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

func (s *GreeterService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	return s.uc.Login(ctx, req)
}

func (s *GreeterService) UploadAsset(ctx context.Context, req *v1.AssetRequest) (*v1.StatusReply, error) {
	return s.uc.UploadAsset(ctx, req)
}

func (s *GreeterService) Get(ctx context.Context, req *v1.AssetRequest) (*v1.GetReply, error) {
	return s.uc.Get(ctx, req)
}
