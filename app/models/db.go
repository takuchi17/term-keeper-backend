package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/takuchi17/term-keeper/configs"
)

const (
	InstanceMySQL int = iota
	InstanceSqlite
)

var DB *sql.DB

func CreateNewDBConnector(instance int) error {
	var err error
	switch instance {
	case InstanceMySQL:
		driver := "mysql"
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			configs.Config.DBUser,
			configs.Config.DBPassword,
			configs.Config.DBHost,
			configs.Config.DBPort,
			configs.Config.DBName,
		)
		err = setupDatabase(driver, dsn)

	case InstanceSqlite:
		driver := "sqlite3"
		dsn := "./test.sql"
		err = setupDatabase(driver, dsn)

	default:
		return errors.New("invalid sql db instance")
	}

	if err != nil {
		slog.Error("Failed to setup database", "err", err)
		panic(err)
	}
	return nil
}

func setupDatabase(driver, dsn string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	DB = db
	return nil
}
