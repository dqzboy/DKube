package main

import (
	"dkube/config"
	"dkube/controller"
	"dkube/service"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化k8s client
	service.K8s.Init()
	//初始化Gin对象
	r := gin.Default()
	controller.Router.InitApiRouter(r)
	//启动gin程序
	r.Run(config.ListenAddr)
}
