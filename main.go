package main

import (
	"dkube/config"
	"dkube/controller"
	"dkube/db"
	"dkube/middle"
	"dkube/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	db.Init()
	service.K8s.Init()
	r := gin.Default()
	r.Use(middle.Cors())
	controller.Router.InitApiRouter(r)

	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(":8081", nil)
	}()

	r.Run(config.ListenAddr)
	db.Close()
}
