package commonUtils

import (
	"fmt"
	"github.com/lijie-keith/go_init_project/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	localhost = "127.0.0.1"
	port      = 3306
	dbname    = "go_test"
	username  = "root"
	password  = "123456"
)

var DB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, localhost, port, dbname)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName:                    "",
		ServerVersion:                 "",
		DSN:                           dsn, // 连接数据库信息
		DSNConfig:                     nil,
		Conn:                          nil,
		SkipInitializeWithVersion:     false,
		DefaultStringSize:             0,
		DefaultDatetimePrecision:      nil,
		DisableWithReturning:          false,
		DisableDatetimePrecision:      false,
		DontSupportRenameIndex:        false,
		DontSupportRenameColumn:       false,
		DontSupportForShareClause:     false,
		DontSupportNullAsDefaultValue: false,
		DontSupportRenameColumnUnique: false,
	}), &gorm.Config{ //&gorm.Config 后面的参数是相关配置，可以根据开发进行修改
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   nil,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		TranslateError:                           false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})
	// 判断连接情况
	if err != nil {
		config.SystemLogger.Info("----------数据库连接失败--------------")
		//panic 抛出异常，并终止程序
		panic(err)
	}
	config.SystemLogger.Info("----------数据库连接成功--------------")
	// 分别赋值给你全局变量DB 和DBD
	DB = db
}
