package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Node node

type node struct{}

type NodesResp struct {
	Items []corev1.Node `json:"items"`
	Total int           `json:"total"`
}

func (n *node) GetNodes(filterName string, limit, page int) (nodesResp *NodesResp, err error) {
	nodeList, err := K8s.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Node列表失败, " + err.Error()))
		return nil, errors.New("获取Node列表失败, " + err.Error())
	}
	selectableData := &dataSelector{
		GenericDataList: n.toCells(nodeList.Items),
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

	nodes := n.fromCells(data.GenericDataList)

	return &NodesResp{
		Items: nodes,
		Total: total,
	}, nil
}

func (n *node) GetNodeDetail(nodeName string) (node *corev1.Node, err error) {
	node, err = K8s.ClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Node详情失败, " + err.Error()))
		return nil, errors.New("获取Node详情失败, " + err.Error())
	}

	return node, nil
}

func (n *node) toCells(std []corev1.Node) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = nodeCell(std[i])
	}
	return cells
}

func (n *node) fromCells(cells []DataCell) []corev1.Node {
	nodes := make([]corev1.Node, len(cells))
	for i := range cells {
		nodes[i] = corev1.Node(cells[i].(nodeCell))
	}

	return nodes
}
