package controller

import (
	"dkube/service"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
)

var Node node

type node struct{}

func (n *node) GetNodes(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
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

	data, err := service.Node.GetNodes(params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Node列表成功",
		"data": data,
	})
}

func (n *node) GetNodeDetail(ctx *gin.Context) {
	params := new(struct {
		NodeName string `form:"node_name"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Node.GetNodeDetail(params.NodeName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Node详情成功",
		"data": data,
	})
}
