package models

import (
	"context"
	"testing"

	"github.com/takuchi17/term-keeper/pkg/tester"
)

func SetupTestDB(t *testing.T, ctx context.Context) (*tester.MysqlContainer, error) {
	container, err := tester.SetupMySQL(ctx)
	if err != nil {
		return nil, err
	}
	// グローバル変数への代入は呼び出し元で実施する
	return container, nil
}
