package config

import "time"

const (
	ListenAddr = "0.0.0.0:9090"
	Kubeconfig = "/root/.kube/config"

	//tail的日志行数
	PodLogTailLine = 2000

	//数据库配置
	DbType = "mysql"
	DbHost = "192.168.66.62"
	DbPort = 3306
	DbName = "dkube"
	DbUser = "root"
	DbPwd  = "123456"

	LogMode = false

	MaxIdleConns = 10
	MaxOpenConns = 100
	MaxLifeTime  = 30 * time.Second

	AdminUser = "admin"
	AdminPwd  = "123456"
)
