package models

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/takuchi17/term-keeper/pkg/tester"
)

var mysqlContainer *tester.MysqlContainer
var DB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	mysqlContainer, err = SetupTestDB(nil, ctx)
	if err != nil {
		log.Fatalf("Failed to setup MySQL container: %v", err)
	}

	// グローバルな DB 変数に代入
	DB, err = mysqlContainer.OpenDB(ctx)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	code := m.Run()

	if err := mysqlContainer.Terminate(ctx); err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}

	os.Exit(code)
}
