package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var StatefulSet statefulSet

type statefulSet struct{}

type StatusfulSetsResp struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

func (s *statefulSet) GetStatefulSets(filterName, namespace string, limit, page int) (statusfulSetsResp *StatusfulSetsResp, err error) {
	statefulSetList, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取StatefulSet列表失败, " + err.Error()))
		return nil, errors.New("获取StatefulSet列表失败, " + err.Error())
	}
	selectableData := &dataSelector{
		GenericDataList: s.toCells(statefulSetList.Items),
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
	statefulSets := s.fromCells(data.GenericDataList)

	return &StatusfulSetsResp{
		Items: statefulSets,
		Total: total,
	}, nil
}

func (s *statefulSet) GetStatefulSetDetail(statefulSetName, namespace string) (statefulSet *appsv1.StatefulSet, err error) {
	statefulSet, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取StatefulSet详情失败, " + err.Error()))
		return nil, errors.New("获取StatefulSet详情失败, " + err.Error())
	}

	return statefulSet, nil
}

func (s *statefulSet) DeleteStatefulSet(statefulSetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除StatefulSet失败, " + err.Error()))
		return errors.New("删除StatefulSet失败, " + err.Error())
	}

	return nil
}

func (s *statefulSet) UpdateStatefulSet(namespace, content string) (err error) {
	var statefulSet = &appsv1.StatefulSet{}

	err = json.Unmarshal([]byte(content), statefulSet)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新StatefulSet失败, " + err.Error()))
		return errors.New("更新StatefulSet失败, " + err.Error())
	}
	return nil
}

func (s *statefulSet) toCells(std []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = statefulSetCell(std[i])
	}
	return cells
}

func (s *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulSets := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		statefulSets[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}

	return statefulSets
}
