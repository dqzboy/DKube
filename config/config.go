package config

import "time"

const (
	ListenAddr = "0.0.0.0:9090"
	//k8s集群kube-config文件存放位置
	Kubeconfig = "H:\\project\\kubeconfig\\config"
	//tail的日志行数
	PodLogTailLine = 2000

	//数据库配置
	DbType = "mysql"
	DbHost = "192.168.66.62"
	DbPort = 3306
	DbName = "dkube"
	DbUser = "root"
	DbPwd  = "123456"

	//打印Mysql debug的sql日志
	LogMode = false

	//连接池的配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
)
