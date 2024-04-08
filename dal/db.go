package dal

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	mysqlDSN    = "root:root1234@tcp(127.0.0.1:13306)/db3?charset=utf8mb4&parseTime=True"
	sqliteTable = "./test.db"
)

func ConnectMySQL(datasource ...string) *gorm.DB {
	var dsn string
	if len(datasource) == 0 {
		dsn = mysqlDSN
	} else {
		dsn = datasource[0]
	}
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	return db
}

func ConnectSQLite(tablename ...string) *gorm.DB {
	var table string
	if len(tablename) == 0 {
		table = sqliteTable
	} else {
		table = tablename[0]
	}
	db, err := gorm.Open(sqlite.Open(table))
	if err != nil {
		panic(fmt.Errorf("sqlite.Open(table) failed, err: %v", err))
	}
	return db
}
