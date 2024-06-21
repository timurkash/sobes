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

func (s *GreeterService) CreateRoute(ctx context.Context, req *v1.CreateRouteRequest) (*v1.CreateRouteReply, error) {
	return s.uc.CreateRouter(ctx, req)
}

func (s *GreeterService) GetRoute(ctx context.Context, req *v1.GetRouteRequest) (*v1.RouteReply, error) {
	return s.uc.GetRouter(ctx, req)
}

func (s *GreeterService) DeleteRoute(ctx context.Context, req *v1.DeleteRouteRequest) (*v1.Empty, error) {
	return s.uc.DeleteRoute(ctx, req)
}
