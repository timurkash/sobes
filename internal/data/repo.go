package data

import (
	"context"
	v1 "helloworld/api/helloworld/v1"
	"helloworld/internal/biz"
	"helloworld/internal/conf"
)

type greeterRepo struct {
	data           *Data
	availableRooms []*conf.Room
	orders         []*v1.CreateOrderRequest
}

// NewGreeterRepo .
func NewGreeterRepo(confData *conf.Data, data *Data) (biz.GreeterRepo, func(), error) {
	return &greeterRepo{
			data:           data,
			availableRooms: confData.RoomAvailability.Rooms,
		}, func() {
			data.log.Info("closing the data resources")
		}, nil
}

func (r *greeterRepo) ListAvailableRooms(_ context.Context) ([]*conf.Room, error) {
	return r.availableRooms, nil
}

func (r *greeterRepo) CreateOrder(_ context.Context, order *v1.CreateOrderRequest) error {
	r.orders = append(r.orders, order)
	return nil
}
