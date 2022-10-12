package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	nwv1 "k8s.io/api/networking/v1"
	"sort"
	"strings"
	"time"
)

type dataSelector struct {
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

type DataSelectQuery struct {
	Filter   *FilterQuery
	Paginate *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()

	return b.Before(a)
}

func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

func (d *dataSelector) Filter() *dataSelector {
	if d.DataSelect.Filter.Name == "" {
		return d
	}
	filtered := []DataCell{}
	for _, value := range d.GenericDataList {
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelect.Filter.Name) {
			matches = false
			continue
		}
		if matches {
			filtered = append(filtered, value)
		}
	}
	d.GenericDataList = filtered
	return d
}

func (d *dataSelector) Paginate() *dataSelector {
	limit := d.DataSelect.Paginate.Limit
	page := d.DataSelect.Paginate.Page
	if limit <= 0 || page <= 0 {
		return d
	}
	startIndex := limit * (page - 1)
	endIndex := limit*page - 1

	if endIndex > len(d.GenericDataList) {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]

	return d
}

type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentCell) GetName() string {
	return d.Name
}

type daemonSetCell appsv1.DaemonSet

func (d daemonSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d daemonSetCell) GetName() string {
	return d.Name
}

type statefulSetCell appsv1.StatefulSet

func (s statefulSetCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s statefulSetCell) GetName() string {
	return s.Name
}

type nodeCell corev1.Node

func (n nodeCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n nodeCell) GetName() string {
	return n.Name
}

type namespaceCell corev1.Namespace

func (n namespaceCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n namespaceCell) GetName() string {
	return n.Name
}

type pvCell corev1.PersistentVolume

func (p pvCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvCell) GetName() string {
	return p.Name
}

type serviceCell corev1.Service

func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceCell) GetName() string {
	return s.Name
}

type ingressCell nwv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressCell) GetName() string {
	return i.Name
}

type configMapCell corev1.ConfigMap

func (c configMapCell) GetCreation() time.Time {
	return c.CreationTimestamp.Time
}

func (c configMapCell) GetName() string {
	return c.Name
}

type secretCell corev1.Secret

func (s secretCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s secretCell) GetName() string {
	return s.Name
}

type pvcCell corev1.PersistentVolumeClaim

func (p pvcCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvcCell) GetName() string {
	return p.Name
}
