package tester

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/docker/go-connections/nat"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MysqlContainer struct {
	testcontainers.Container
}

func (s *MysqlContainer) OpenDB(ctx context.Context) (*sql.DB, error) {
	host, err := s.Container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := s.Container.MappedPort(ctx, "3306")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"user",
		"password",
		host,
		port.Port(),
		"term_keeper_db_test",
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupMySQL(ctx context.Context) (*MysqlContainer, error) {
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
			"MYSQL_CHARACTER_SET_SERVER": "utf8mb4",
			"MYSQL_COLLATION_SERVER":     "utf8mb4_unicode_ci",
		},
		ExposedPorts: []string{"3306/tcp"},
		Mounts: testcontainers.ContainerMounts{
			testcontainers.BindMount(initdbDir, "/docker-entrypoint-initdb.d"),
		},
		WaitingFor: wait.ForSQL("3306", "mysql", func(host string, port nat.Port) string {
			return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				"user",
				"password",
				host,
				port.Port(),
				"term_keeper_db_test",
			)
		}),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &MysqlContainer{
		Container: container,
	}, nil
}
