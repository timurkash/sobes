package data

import (
	"context"
	"helloworld/internal/conf"
	"helloworld/internal/data/ent"
	"helloworld/internal/data/ent/migrate"

	"helloworld/internal/biz"
)

type greeterRepo struct {
	data       *Data
	relational *ent.Client
}

// NewGreeterRepo .
func NewGreeterRepo(confData *conf.Data, data *Data) (biz.GreeterRepo, func(), error) {
	driver, err := GetRelationalDriver(confData.Relational)
	if err != nil {
		return nil, nil, err
	}
	entClient := ent.NewClient(ent.Driver(DebugWithContext(driver, data.log)))
	if err := entClient.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		data.log.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}
	return &greeterRepo{
			data: data,
		}, func() {
			data.log.Info("closing the data resources")
			if err := entClient.Close(); err != nil {
				data.log.Error(err)
			}
		}, nil
}

func (r *greeterRepo) Create(ctx context.Context, route *biz.Route) error {
	_, err := r.relational.Route.Create().
		SetID(route.Id).
		SetRouteName(route.Name).
		SetLoad(route.Load).
		SetCargoType(route.CargoType).
		SetIsActual(route.IsActual).
		Save(ctx)
	return err
}

func (r *greeterRepo) FindByID(ctx context.Context, id uint64) (*biz.Route, error) {
	record, err := r.relational.Route.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &biz.Route{
		Id:        id,
		Name:      record.RouteName,
		Load:      record.Load,
		CargoType: record.CargoType,
		IsActual:  record.IsActual,
	}, nil
}

func (r *greeterRepo) SetIsActual(ctx context.Context, id uint64, isActual bool) error {
	_, err := r.relational.Route.UpdateOneID(id).SetIsActual(isActual).Save(ctx)
	return err
}

func (r *greeterRepo) DeleteById(ctx context.Context, id uint64) error {
	return r.relational.Route.DeleteOneID(id).Exec(ctx)
}
