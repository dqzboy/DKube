package controller

import (
	"dkube/service"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
)

var ConfigMap configMap

type configMap struct{}

func (c *configMap) GetConfigMaps(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.ConfigMap.GetConfigMaps(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取ConfigMap列表成功",
		"data": data,
	})
}

func (c *configMap) GetConfigMapDetail(ctx *gin.Context) {
	params := new(struct {
		ConfigMapName string `form:"configmap_name"`
		Namespace     string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.ConfigMap.GetConfigMapDetail(params.ConfigMapName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取ConfigMap详情成功",
		"data": data,
	})
}

func (c *configMap) DeleteConfigMap(ctx *gin.Context) {
	params := new(struct {
		ConfigMapName string `json:"configmap_name"`
		Namespace     string `json:"namespace"`
	})

	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.ConfigMap.DeleteConfigMap(params.ConfigMapName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除ConfigMap成功",
		"data": nil,
	})
}

func (c *configMap) UpdateConfigMap(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})

	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.ConfigMap.UpdateConfigMap(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新ConfigMap成功",
		"data": nil,
	})
}
