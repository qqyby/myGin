package mysql

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var db *sqlx.DB

func Init(user, pwd, dbname, host string, port, maxOpenConn, maxIdleConn int) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", user, pwd, host, port, dbname)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return errors.Wrapf(err, "mysql connect failed")
	}

	if err := db.Ping(); err != nil {
		return errors.Wrapf(err, "ping db error")
	}

	// 根据具体业务设置
	db.SetMaxOpenConns(maxOpenConn) // 最大连接数
	db.SetMaxIdleConns(maxIdleConn) // 最大空闲连接数
	return nil
}

func Close() error {
	return db.Close()
}

func GetPageOffset(pageNumber, pageSize int) int {
	result := 0
	if pageNumber > 0 {
		result = (pageNumber - 1) * pageSize
	}

	return result
}

func GenerateWhereQuery(wheres []string) string {
	if len(wheres) == 0 {
		return ""
	}

	if len(wheres) == 1 {
		return " WHERE " + wheres[0]
	}

	return " WHERE " + strings.Join(wheres, " AND ")
}

