package tester

import (
	"context"
	"database/sql"
	"net"
	"path/filepath"

	"github.com/docker/go-connections/nat"
	"github.com/go-sql-driver/mysql"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type mysqlContainer struct {
	testcontainers.Container
}

func (s *mysqlContainer) OpenDB(ctx context.Context) (*sql.DB, error) {
	host, err := s.Container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := s.Container.MappedPort(ctx, "3306")
	if err != nil {
		return nil, err
	}

	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(host, port.Port())
	cfg.DBName = "term_keeper_db_test"
	cfg.User = "user"
	cfg.Passwd = "password"

	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, err
	}

	return sql.OpenDB(connector), nil
}

func SetupMySQL(ctx context.Context) (*mysqlContainer, error) {
	initdbDir, err := filepath.Abs("../../pkg/tester/testdata")
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		Image: "mysql:8.0",
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD":        "password",
			"MYSQL_DATABASE":             "term_keeper_db_test",
			"MYSQL_USER":                 "user",
			"MYSQL_PASSWORD":             "password",
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
		},
		ExposedPorts: []string{"3306/tcp"},
		Mounts: testcontainers.ContainerMounts{
			testcontainers.BindMount(initdbDir, "/docker-entrypoint-initdb.d"),
		},
		WaitingFor: wait.ForSQL("3306", "mysql", func(host string, port nat.Port) string {
			cfg := mysql.NewConfig()
			cfg.Net = "tcp"
			cfg.Addr = net.JoinHostPort(host, port.Port())
			cfg.DBName = "term_keeper_db_test"
			cfg.User = "user"
			cfg.Passwd = "password"
			return cfg.FormatDSN()
		}),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &mysqlContainer{
		Container: container,
	}, nil
}
