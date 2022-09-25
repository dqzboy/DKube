package controller

import (
	"github.com/gin-gonic/gin"
)

var Router router

//声明router结构体
type router struct{}

// InitApiRouter 初始化路由规则，创建api测试接口
func (r *router) InitApiRouter(router *gin.Engine) {
	router.
		//pod操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		GET("/api/k8s/pod/numnp", Pod.GetPodNumPerNp)
}
