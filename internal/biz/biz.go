package biz

import (
	"context"
	"github.com/google/wire"
	"helloworld/internal/data/ent"
	"math/rand"

	v1 "helloworld/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase)

var (
	ErrRouteAlreadyExists = errors.New(208, v1.ErrorReason_ROUTE_ALREADY_EXISTS.String(), "")
	ErrRouteNotActual     = errors.New(410, v1.ErrorReason_ROUTE_NOT_ACTUAL.String(), "")
	ErrRouteInProcess     = errors.New(202, v1.ErrorReason_ROUTE_IN_PROCESS.String(), "")
)

type Route struct {
	Id        uint64
	Name      string
	Load      float64
	CargoType string
	IsActual  bool
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Create(ctx context.Context, router *Route) error
	SetIsActual(ctx context.Context, id uint64, isActual bool) error
	FindByID(ctx context.Context, id uint64) (*Route, error)
	DeleteById(ctx context.Context, uid uint64) error
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

func (uc *GreeterUsecase) CreateRouter(ctx context.Context, req *v1.CreateRouteRequest) (*v1.CreateRouteReply, error) {
	// Проверка входных данных
	found := true
	_, err := uc.repo.FindByID(ctx, req.RouteId)
	if err != nil {
		if ent.IsNotFound(err) {
			found = false
		} else {
			return nil, err
		}
	}
	if found {
		if err := uc.repo.SetIsActual(ctx, req.RouteId, false); err != nil {
			return nil, err
		}
		req.RouteId = rand.Uint64()
	}
	if err := uc.repo.Create(ctx, &Route{
		Id:        req.RouteId,
		Name:      req.RouteName,
		Load:      req.Load,
		CargoType: req.CargoType,
		IsActual:  true,
	}); err != nil {
		return nil, err
	}
	return &v1.CreateRouteReply{
		RouteId:       req.RouteId,
		RouteName:     req.RouteName,
		Load:          req.Load,
		CargoType:     req.CargoType,
		AlreadyExists: found,
	}, nil
}

func (uc *GreeterUsecase) GetRouter(ctx context.Context, req *v1.GetRouteRequest) (*v1.RouteReply, error) {
	route, err := uc.repo.FindByID(ctx, req.RouteId)
	if err != nil {
		return nil, err
	}
	if !route.IsActual {
		return nil, ErrRouteNotActual
	}
	return &v1.RouteReply{
		RouteName: route.Name,
		Load:      route.Load,
		CargoType: route.CargoType,
	}, nil
}

func (uc *GreeterUsecase) DeleteRoute(ctx context.Context, req *v1.DeleteRouteRequest) (*v1.Empty, error) {
	if err := uc.repo.DeleteById(ctx, req.RouteId); err != nil {
		return nil, err
	}
	return &v1.Empty{}, nil
}
