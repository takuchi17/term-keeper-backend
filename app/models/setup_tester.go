package models

import (
	"context"
	"testing"

	"github.com/takuchi17/term-keeper/pkg/tester"
)

func SetupMysqlContainerAndSetupDB(t *testing.T, ctx *context.Context) (*tester.MysqlContainer, error) {
	container, err := tester.SetupMySQL(*ctx)
	DB, err = container.OpenDB(*ctx)
	if err != nil {
		return nil, err
	}
	return container, nil
}
