package service

import (
	"dkube/config"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

func (k *k8s) Init() {
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		logger.Error("获取k8s client配置失败," + err.Error())
	}
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		logger.Error("创建k8s client失败," + err.Error())
	} else {
		logger.Info("k8s client 初始化成功!")
	}

	k.ClientSet = clientset
}
