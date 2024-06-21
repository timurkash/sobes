package data

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	dbSql "database/sql"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	entSql "entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"helloworld/internal/conf"
	"strings"

	// init drivers
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	postgresConnectStringConst = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s %s"
	mysqlConnectStringConst    = "%s:%s@tcp(%s:%d)/%s?parseTime=true"
)

var NoPostgresMySQL = errors.New("no postgres or mysql is not supported")

func GetBadConfigError(config string) error {
	return fmt.Errorf("config %s is wrong", config)
}

func ReplaceEndOfLine(str string) string {
	return strings.ReplaceAll(str, `\n`, "\n")
}

func DebugWithContext(driver *sql.Driver, logHelper *log.Helper) dialect.Driver {
	return dialect.DebugWithContext(driver, func(ctx context.Context, i ...interface{}) {
		logHelper.WithContext(ctx).Info(i...)
		_, span := otel.Tracer("entgo.io").Start(ctx,
			"Query",
			trace.WithAttributes(
				attribute.String("sql", fmt.Sprint(i...)),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		span.End()
	})
}

func SetSchema(driver *sql.Driver, schema string) error {
	if schema != "" && schema != "public" {
		if _, err := driver.DB().Exec(
			fmt.Sprintf(
				"create schema if not exists %s;set search_path to %s;",
				schema,
				schema,
			),
		); err != nil {
			return err
		}
	}
	return nil
}

func GetRelationalDriver(relational *conf.Relational) (*entSql.Driver, error) {
	if relational == nil {
		return nil, GetBadConfigError("relational")
	}
	if relational.Dialect == "" {
		return nil, GetBadConfigError("relational.dialect")
	}
	if relational.Host == "" {
		return nil, GetBadConfigError("relational.host")
	}
	if relational.Port == 0 {
		return nil, GetBadConfigError("relational.port")
	}
	if relational.Dbname == "" {
		return nil, GetBadConfigError("relational.dbname")
	}
	if relational.User == "" {
		return nil, GetBadConfigError("relational.user")
	}
	if relational.Password == "" {
		return nil, GetBadConfigError("relational.password")
	}
	switch relational.Dialect {
	case dialect.Postgres:
		if relational.SslMode == "" {
			return nil, GetBadConfigError("relational.ssl_mode")
		}
		connString := strings.Trim(fmt.Sprintf(postgresConnectStringConst,
			relational.Host,
			relational.Port,
			relational.User,
			relational.Password,
			relational.Dbname,
			relational.SslMode,
			relational.Additional,
		), " ")
		connConfig, err := pgx.ParseConfig(connString)
		if err != nil {
			return nil, err
		}
		if relational.CaCertificate != "" {
			rootCertPool := x509.NewCertPool()
			relational.CaCertificate = ReplaceEndOfLine(relational.CaCertificate)
			if ok := rootCertPool.AppendCertsFromPEM([]byte(relational.CaCertificate)); !ok {
				return nil, errors.New("failed to append PEM")
			}
			connConfig.TLSConfig = &tls.Config{
				RootCAs: rootCertPool,
				//InsecureSkipVerify: true,
			}
		}
		db, err := dbSql.Open("pgx", connString)
		if err != nil {
			return nil, err
		}
		driver := entSql.OpenDB(dialect.Postgres, db)
		if err := SetSchema(driver, relational.Schema); err != nil {
			return nil, err
		}
		return driver, err
	case dialect.MySQL:
		const custom = "custom"
		mysqlConnectString := mysqlConnectStringConst
		if relational.CaCertificate != "" {
			rootCertPool := x509.NewCertPool()
			caCertificate := strings.ReplaceAll(relational.CaCertificate, `\n`, "\n")
			if ok := rootCertPool.AppendCertsFromPEM([]byte(caCertificate)); !ok {
				return nil, errors.New("failed to append PEM")
			}
			if err := mysql.RegisterTLSConfig(custom, &tls.Config{
				RootCAs: rootCertPool,
			}); err != nil {
				return nil, err
			}
			mysqlConnectString = fmt.Sprintf("%s&tls=%s", mysqlConnectStringConst, custom)
		}
		return entSql.Open(
			dialect.MySQL,
			fmt.Sprintf(mysqlConnectString,
				relational.User,
				relational.Password,
				relational.Host,
				relational.Port,
				relational.Dbname,
			),
		)
	}
	return nil, NoPostgresMySQL
}
