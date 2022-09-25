package service

import (
	appsv1 "k8s.io/api/apps/v1" //引入K8s的APPSv1的包.定义别名为appsv1
	corev1 "k8s.io/api/core/v1" //引入K8s的corev1的包.定义别名为corev1
	nwv1 "k8s.io/api/networking/v1"
	"sort"
	"strings"
	"time"
)

//定义结构体：dataSelector，用于封装排序、过滤、分页的数据类型
type dataSelector struct { //包内调用，首字母小写
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

// DataCell DataCell接口用于各种资源List类型的转换，转换后可以使用dataSelector的排序、过滤、分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页结构体,过滤使用Name过滤，分页使用Limit和page
//Limit 是单页的数据条数
//Page是第几页
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

// Len 实现自定义结构体的排序，需要重写Len、Swap、Less方法
//Len方法用于获取数组的长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap 方法用于数据比较大小之后的位置变更
func (d *dataSelector) Swap(i, j int) {
	//临时变量对调  i,j ==> j.i
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 方法用于比较大小
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	//比较B的时间是否在A之前，触发位置调换
	return b.Before(a)
}

//重写以上三个方法，使用sort.Sort进行排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// Filter 方法用于过滤数据，比较数据的Name属性，如果包含则返回
func (d *dataSelector) Filter() *dataSelector {
	//判断入参是否为空，如果为空则返回所有数据
	if d.DataSelect.Filter.Name == "" {
		return d
	}
	//如果不为空，则按照入参Name进行过滤
	//声明一个新的数组。若Name包含，则把数据放进数组，返回出去
	filtered := []DataCell{}
	for _, value := range d.GenericDataList {
		//定义是否匹配的标签变量，然后默认是匹配
		matches := true
		objName := value.GetName()
		if strings.Contains(objName, d.DataSelect.Filter.Name) {
			matches = false
			continue //跳过当前循环，执行下一次循环
		}
		if matches {
			filtered = append(filtered, value)
		}
	}

	d.GenericDataList = filtered
	return d
}

// Paginate 方法用于数组的分页，根据Limit和Page的传参，取一定范围内的数据，返回
func (d *dataSelector) Paginate() *dataSelector {
	//根据Limit和Page的入参，定义快捷变量
	limit := d.DataSelect.Paginate.Limit
	page := d.DataSelect.Paginate.Page

	//检验参数的合法性
	if limit <= 0 || page <= 0 {
		return d
	}

	//定义取数范围需要的starIndex和endIndex
	startIndex := limit * (page - 1)
	endIndex := limit*page - 1

	//出来endIndex
	if endIndex > len(d.GenericDataList) {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]

	return d
}

// 定义podCell类型，重写GetCreateion和GetName方法后，可进行数据转换
//corev1.Pod --> podCell --> DataCell
//appsv1.Deployment --> deployCell --> DataCell
type podCell corev1.Pod

// GetCreation 重写DataCell接口的两个方法
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

type serviceCell corev1.Service

func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceCell) GetName() string {
	return s.Name
}

//追加ingress DataCell接口相关代码
type ingressCell nwv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressCell) GetName() string {
	return i.Name
}
