package models

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takuchi17/term-keeper/configs"
)

const InstanceMySQL int = iota

type SQLExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func NewDB(instance int) (*sql.DB, error) {
	switch instance {
	case InstanceMySQL:
		driver := "mysql"
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			configs.Config.DBUser,
			configs.Config.DBPassword,
			configs.Config.DBHost,
			configs.Config.DBPort,
			configs.Config.DBName,
		)
		db, err := sql.Open(driver, dsn)
		if err != nil {
			return nil, err
		}
		return db, nil

	default:
		return nil, errors.New("invalid sql db instance")
	}
}
