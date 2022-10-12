package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pv pv

type pv struct{}

type PvsResp struct {
	Items []corev1.PersistentVolume `json:"items"`
	Total int                       `json:"total"`
}

func (p *pv) GetPvs(filterName string, limit, page int) (pvsResp *PvsResp, err error) {
	pvList, err := K8s.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Pv列表失败, " + err.Error()))
		return nil, errors.New("获取Pv列表失败, " + err.Error())
	}
	selectableData := &dataSelector{
		GenericDataList: p.toCells(pvList.Items),
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

	pvs := p.fromCells(data.GenericDataList)

	return &PvsResp{
		Items: pvs,
		Total: total,
	}, nil
}

func (p *pv) GetPvDetail(pvName string) (pv *corev1.PersistentVolume, err error) {
	pv, err = K8s.ClientSet.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Pv详情失败, " + err.Error()))
		return nil, errors.New("获取Pv详情失败, " + err.Error())
	}

	return pv, nil
}

func (p *pv) DeletePv(pvName string) (err error) {
	err = K8s.ClientSet.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Pv失败, " + err.Error()))
		return errors.New("删除Pv失败, " + err.Error())
	}

	return nil
}

func (p *pv) toCells(std []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = pvCell(std[i])
	}
	return cells
}

func (p *pv) fromCells(cells []DataCell) []corev1.PersistentVolume {
	pvs := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		pvs[i] = corev1.PersistentVolume(cells[i].(pvCell))
	}

	return pvs
}
