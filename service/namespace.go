package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Namespace namespace

type namespace struct{}

type NamespacesResp struct {
	Items []corev1.Namespace `json:"items"`
	Total int                `json:"total"`
}

func (n *namespace) GetNamespaces(filterName string, limit, page int) (namespacesResp *NamespacesResp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace列表失败, " + err.Error()))
		return nil, errors.New("获取Namespace列表失败, " + err.Error())
	}
	selectableData := &dataSelector{
		GenericDataList: n.toCells(namespaceList.Items),
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
	data := filtered.Sort().Paginate()

	namespaces := n.fromCells(data.GenericDataList)

	return &NamespacesResp{
		Items: namespaces,
		Total: total,
	}, nil
}

func (n *namespace) GetNamespaceDetail(namespaceName string) (namespace *corev1.Namespace, err error) {
	namespace, err = K8s.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace详情失败, " + err.Error()))
		return nil, errors.New("获取Namespace详情失败, " + err.Error())
	}

	return namespace, nil
}

func (n *namespace) DeleteNamespace(namespaceName string) (err error) {
	err = K8s.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Namespace失败, " + err.Error()))
		return errors.New("删除Namespace失败, " + err.Error())
	}

	return nil
}

func (n *namespace) toCells(std []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = namespaceCell(std[i])
	}
	return cells
}

func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	namespaces := make([]corev1.Namespace, len(cells))
	for i := range cells {
		namespaces[i] = corev1.Namespace(cells[i].(namespaceCell))
	}

	return namespaces
}
