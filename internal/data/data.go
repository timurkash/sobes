package data

import (
	"helloworld/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	log *log.Helper
}

// NewData .
func NewData(confData *conf.Data, logger log.Logger) (*Data, error) {
	if confData == nil {
		return nil, GetBadConfigError("data")
	}
	return &Data{
		log: log.NewHelper(logger),
	}, nil
}
