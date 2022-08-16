package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Pod pod

type pod struct{}

//Controller中的方法入参是gin.Context，用于从上下文中获取请求参数以及定义响应内容
//流程：绑定参数 --> 调用service代码 --> 根据调用结果响应具体内容

// GetPods 获取Pod列表，支持分页、过滤、排序
func (p *pod) GetPods(ctx *gin.Context) {
	//处理入参
	//匿名结构体用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	//form格式使用Bind方法，json格式使用ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind参数绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": data,
	})
}

//获取Pod详情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	//处理入参
	//匿名结构体用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	//form格式使用Bind方法，json格式使用ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind参数绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": data,
	})
}

//删除Pod
func (p *pod) DeletePod(ctx *gin.Context) {
	//处理入参
	//匿名结构体用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})
	//form格式使用Bind方法，json格式使用ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind参数绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Pod成功",
		"data": nil,
	})
}

//更新Pod
func (p *pod) UpdatePod(ctx *gin.Context) {
	//处理入参
	//匿名结构体用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})
	//PUT请求，绑定参数方法改为ctx.ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind参数绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Pod.UpdatePod(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新Pod成功",
		"data": nil,
	})
}

//获取Pod容器
func (p *pod) GetPodContainer(ctx *gin.Context) {
	//处理入参
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	//GET请求，绑定参数方法改为ctx.Bind方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPodContainer(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod容器成功",
		"data": data,
	})
}

//获取Pod容器日志
func (p *pod) GetPodLog(ctx *gin.Context) {
	//处理入参
	params := new(struct {
		ContainerName string `form:"container_name"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
	})
	//GET请求，绑定参数方法改为ctx.Bind方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod容器日志成功",
		"data": data,
	})
}

//获取每个Namespace的Pod数量
func (p *pod) GetPodNumPerNp(ctx *gin.Context) {
	data, err := service.Pod.GetPodNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取每个namespace的Pod数量成功",
		"data": data,
	})
}
