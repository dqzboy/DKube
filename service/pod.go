package service

import (
	"bytes"
	"context"
	"dkube/config"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct{}

type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}


type PodsNp struct {
	Namespace string
	PodNum    int
}


func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取Pod列表失败," + err.Error())
		return nil, errors.New("获取Pod列表失败," + err.Error())
	}
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

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	//排序和分页
	data := filtered.Sort().Paginate()

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
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

//获取Pod内容器日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	//获取request实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(errors.New("获取PodLogs失败," + err.Error()))
		return "", errors.New("获取PodLogs失败," + err.Error())
	}
	defer podLogs.Close() //处理完关闭
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
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum:    len(podList.Items),
		}
		podsNps = append(podsNps, podsNp)
	}
	return podsNps, nil
}

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
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
