package data

import (
	"context"
	"helloworld/internal/conf"
	"helloworld/internal/data/ent"
	"helloworld/internal/data/ent/session"
	"helloworld/internal/data/ent/user"

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
	//if err := entClient.Schema.Create(
	//	context.Background(),
	//	migrate.WithDropIndex(true),
	//	migrate.WithDropColumn(true),
	//); err != nil {
	//	data.log.Errorf("failed creating schema resources: %v", err)
	//	return nil, nil, err
	//}
	return &greeterRepo{
			data:       data,
			relational: entClient,
		}, func() {
			data.log.Info("closing the data resources")
			if err := entClient.Close(); err != nil {
				data.log.Error(err)
			}
		}, nil
}

func (r *greeterRepo) GetPasswordHash(ctx context.Context, login string) (uint64, string, error) {
	record, err := r.relational.User.Query().Where(user.Login(login)).Only(ctx)
	if err != nil {
		return 0, "", err
	}
	return record.ID, record.PasswordHash, nil
}

func (r *greeterRepo) CreateSession(ctx context.Context, userId uint64) (string, error) {
	record, err := r.relational.Session.Create().SetUID(userId).Save(ctx)
	if err != nil {
		return "", err
	}
	return record.ID, nil
}

func (r *greeterRepo) GetUserId(ctx context.Context, token string) (uint64, error) {
	record, err := r.relational.Session.Get(ctx, token)
	if err != nil {
		return 0, nil
	}
	return record.UID, nil
}

func (r *greeterRepo) DeleteAllSession(ctx context.Context, userId uint64) error {
	_, err := r.relational.Session.Delete().Where(session.UID(userId)).Exec(ctx)
	return err
}

func (r *greeterRepo) InsertAsset(ctx context.Context, userId uint64, assetName string, file []byte) error {
	//TODO implement me
	panic("implement me")
}

func (r *greeterRepo) GetAsset(ctx context.Context, assetName string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
