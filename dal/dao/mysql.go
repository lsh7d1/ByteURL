package dao

import (
	"context"
	"fmt"

	"byteurl/dal/model"
	"byteurl/dal/query"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const mysqldsn = "root:root1234@tcp(127.0.0.1:13306)/db2?charset=utf8mb4&parseTime=True"

func init() {
	connectMySQL(mysqldsn)
}

func connectMySQL(dsn string) {
	mdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("mysql.Open(dsn) failed, err: %v", err))
	}
	query.SetDefault(mdb)
	// mdb.AutoMigrate(&model.Short{})
}

// InsertShortURL 插入一条short数据
func InsertShortURL(ctx context.Context, s *model.Short) error {
	return query.Short.WithContext(ctx).Create(s)
}

// PeekMd5NotExist 查看给定的长链md5是否不存在
func PeekMd5NotExist(ctx context.Context, md5 string) error {
	_, err := query.Short.WithContext(ctx).Where(query.Short.Md5.Eq(md5)).First()
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return ErrExistData
}

// PeekShortURLNotExist 查看给定的shorturl是否不存在
func PeekShortURLNotExist(ctx context.Context, surl string) error {
	_, err := query.Short.WithContext(ctx).Where(query.Short.Surl.Eq(surl)).First()
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return ErrExistData
}
