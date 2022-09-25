package main

import (
	"dkube/config"
	"dkube/controller"
	"dkube/service"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化k8s client
	service.K8s.Init() //可以使用service.K8s.ClientSet 跨包调用
	//初始化Gin对象
	r := gin.Default()
	//跨包调用router的方法，初始化路由规则
	controller.Router.InitApiRouter(r)
	//启动gin程序
	r.Run(config.ListenAddr)
}
