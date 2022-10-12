package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ConfigMap configMap

type configMap struct{}

type ConfigMapsResp struct {
	Items []corev1.ConfigMap `json:"items"`
	Total int                `json:"total"`
}

func (c *configMap) GetConfigMaps(filterName, namespace string, limit, page int) (configMapsResp *ConfigMapsResp, err error) {
	configMapList, err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取ConfigMap列表失败, " + err.Error()))
		return nil, errors.New("获取ConfigMap列表失败, " + err.Error())
	}
	selectableData := &dataSelector{
		GenericDataList: c.toCells(configMapList.Items),
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

	configMaps := c.fromCells(data.GenericDataList)

	return &ConfigMapsResp{
		Items: configMaps,
		Total: total,
	}, nil
}

func (c *configMap) GetConfigMapDetail(configMapName, namespace string) (configMap *corev1.ConfigMap, err error) {
	configMap, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取ConfigMap详情失败, " + err.Error()))
		return nil, errors.New("获取ConfigMap详情失败, " + err.Error())
	}

	return configMap, nil
}

func (c *configMap) DeleteConfigMap(configMapName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configMapName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除ConfigMap失败, " + err.Error()))
		return errors.New("删除ConfigMap失败, " + err.Error())
	}

	return nil
}

func (c *configMap) UpdateConfigMap(namespace, content string) (err error) {
	var configMap = &corev1.ConfigMap{}

	err = json.Unmarshal([]byte(content), configMap)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新ConfigMap失败, " + err.Error()))
		return errors.New("更新ConfigMap失败, " + err.Error())
	}
	return nil
}

func (c *configMap) toCells(std []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = configMapCell(std[i])
	}
	return cells
}

func (c *configMap) fromCells(cells []DataCell) []corev1.ConfigMap {
	configMaps := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		configMaps[i] = corev1.ConfigMap(cells[i].(configMapCell))
	}

	return configMaps
}
