package db

import (
	"dkube/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
	"time"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

//db的初始化函数，与数据库建立连接
func Init() {
	//判断是否已经初始化了
	if isInit {
		return
	}
	//组装连接配置；loc是MySQL的时区设置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPwd,
		config.DbHost,
		config.DbPort,
		config.DbName)
	GORM, err = gorm.Open(config.DbType, dsn)
	if err != nil {
		panic("数据库连接失败," + err.Error())
	}

	//打印sql语句
	GORM.LogMode(config.LogMode)

	//开启连接池
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetConnMaxLifetime(time.Duration(config.MaxLifeTime))

	isInit = true
	logger.Info("数据库初始化成功！")
}

//关闭数据库连接
func Close() error {
	fmt.Println("测试db是否关闭")
	return GORM.Close()
}
