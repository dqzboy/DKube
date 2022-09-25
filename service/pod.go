package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	"io"
	"dkube/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct{}

// PodsResp 定义列表的返回内容，Items是Pod的元素列表，Total是元素数量
type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

//声明类型
type PodsNp struct {
	Namespace string
	PodNum    int
}

// GetPods 获取Pod列表，支持过滤、排序、分页
//context.TODD()用于声明一个空的context上下文，用于List方法内设置这个请求的超时(源码)，这里是常用用法
//metav1.ListOptions{}用于过滤List数据，如使用label，field等
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//打印日志,排错
		logger.Info("获取Pod列表失败," + err.Error())
		//返回给上一层，最终返回给前端，前端打印出该error
		return nil, errors.New("获取Pod列表失败," + err.Error())
	}

	//实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	//先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	//排序和分页
	data := filtered.Sort().Paginate()

	//将DataCell类型转换成Pod
	pods := p.fromCells(data.GenericDataList)

	return &PodsResp{
		Total: total,
		Items: pods,
	}, nil

}

//获取Pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Pod详情失败," + err.Error())
		return nil, errors.New("获取Pod详情失败," + err.Error())
	}
	return pod, nil
}

//删除Pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Pod失败," + err.Error())
		return errors.New("删除Pod失败," + err.Error())
	}
	return nil
}

//更新Pod
func (p *pod) UpdatePod(namespace, content string) (err error) {
	pod := &corev1.Pod{}
	//将json反序列化为pod类型
	err = json.Unmarshal([]byte(content), pod)
	if err != nil {
		logger.Error("反序列化失败," + err.Error())
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Pod失败," + err.Error())
		return errors.New("更新Pod失败," + err.Error())
	}
	return nil
}

//获取Pod容器名列表
func (p *pod) GetPodContainer(podName, namespace string) (containers []string, err error) {
	//获取Pod详情
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		return nil, err
	}
	//从Pod对象中拿到容器名
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

//获取Pod内容器日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	//设置日志的配置，容器名、获取内容的配置
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	//获取request实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	//发起request请求，返回一个io.ReadCloser类型（等同于response.body）
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(errors.New("获取PodLogs失败," + err.Error()))
		return "", errors.New("获取PodLogs失败," + err.Error())
	}
	defer podLogs.Close() //处理完关闭
	//将response body写入到缓冲区，目的是为了转成string返回
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error(errors.New("获取PodLogs失败," + err.Error()))
		return "", errors.New("获取PodLogs失败," + err.Error())
	}

	return buf.String(), nil
}

//获取每个namespace的pod数量
func (p *pod) GetPodNumPerNp() (podsNps []*PodsNp, err error) {
	//获取namespace列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		//获取pod列表
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		//组装数据
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum:    len(podList.Items),
		}
		//添加到podsNps数组中
		podsNps = append(podsNps, podsNp)
	}
	return podsNps, nil
}

//类型转换的方法 corev1.Pod -> DataCell, DataCell --> corev1.Pod
func (p *pod) toCells(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//cells[i].(podCell)是将DataCell类型转成podCell类型
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
