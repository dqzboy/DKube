package service

import (
	"dkube/config"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//用于初始化k8s clientset
var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

//初始化方法
func (k *k8s) Init() {
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		panic("获取k8s clinet配置失败," + err.Error())
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("创建k8s clinet失败," + err.Error())
	} else {
		logger.Info("k8s client 初始化成功")
	}

	k.ClientSet = clientset
}
