package database

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
	"time"
	"top-ping/pkg/logger"
)

var (
	DB       *gorm.DB
	initOnce sync.Once
)

func Init(dataSource *DatasourceConfig) {
	initOnce.Do(func() {
		DB = NewMysqlDB(dataSource)
	})
}

func NewMysqlDB(dataSource *DatasourceConfig) *gorm.DB {
	host := dataSource.Addr
	database := dataSource.Database
	username := dataSource.User
	password := dataSource.Password
	charset := dataSource.Charset

	// 拼接mysql相关配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", username, password, host, 3306, database, charset)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	// 打开链接
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.NewGormLogger(),
	})

	if err != nil {
		logger.Errorf(context.Background(), "gorm mysql start failed: %v", err)
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
