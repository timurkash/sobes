package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "helloworld/api/helloworld/v1"
	"helloworld/internal/conf"
	"time"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase)

var (
	BadRequestError       = errors.New(400, "bad request", "")
	NotAuthorizedError    = errors.New(401, "not authorized", "")
	RoomNotAvailableError = errors.New(500, "Hotel room is not available for selected dates", "")

	//ErrRouteNotActual  = errors.New(410, v1.ErrorReason_ROUTE_NOT_ACTUAL.String(), "")
	//ErrRouteInProcess  = errors.New(202, v1.ErrorReason_ROUTE_IN_PROCESS.String(), "")
)

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	ListAvailableRooms(ctx context.Context) ([]*conf.Room, error)
	CreateOrder(ctx context.Context, request *v1.CreateOrderRequest) error
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *kratosLog.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger kratosLog.Logger) *GreeterUsecase {
	return &GreeterUsecase{
		repo: repo,
		log:  kratosLog.NewHelper(logger),
	}
}

func (uc *GreeterUsecase) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.CreateOrderReply, error) {
	daysToBook := daysBetween(req.From.AsTime(), req.To.AsTime())
	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}
	availableRooms, err := uc.repo.ListAvailableRooms(ctx)
	if err != nil {
		return nil, err
	}

	for _, dayToBook := range daysToBook {
		for i, availability := range availableRooms {
			if !availability.Date.AsTime().Equal(dayToBook) || availability.Quota < 1 {
				continue
			}
			availability.Quota -= 1
			availableRooms[i] = availability
			delete(unavailableDays, dayToBook)
		}
	}
	if len(unavailableDays) != 0 {
		return nil, RoomNotAvailableError
	}
	if err := uc.repo.CreateOrder(ctx, req); err != nil {
		return nil, err
	}
	return &v1.CreateOrderReply{}, nil
}
