package data

import (
	"errors"
	"fmt"
)

const (
	postgresConnectStringConst = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s %s"
	mysqlConnectStringConst    = "%s:%s@tcp(%s:%d)/%s?parseTime=true"
)

var NoPostgresMySQL = errors.New("no postgres or mysql is not supported")

func GetBadConfigError(config string) error {
	return fmt.Errorf("config %s is wrong", config)
}
