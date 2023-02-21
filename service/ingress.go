package service
import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct{}
type IngressesResp struct {
	Items []nwv1.Ingress `json:"items"`
	Total int            `json:"total"`
}

type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
}

type HttpPath struct {
	Path        string        `json:"path"`
	PathType    nwv1.PathType `json:"path_type"`
	ServiceName string        `json:"service_name"`
	ServicePort int32         `json:"service_port"`
}

func (i *ingress) GetIngresses(filterName, namespace string, limit, page int) (ingressesResp *IngressesResp, err error) {
	ingressList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress列表失败, " + err.Error()))
		return nil, errors.New("获取Ingress列表失败, " + err.Error())
	}
	selectableData := &dataSelector{
		GenericDataList: i.toCells(ingressList.Items),
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

	ingresss := i.fromCells(data.GenericDataList)

	return &IngressesResp{
		Items: ingresss,
		Total: total,
	}, nil
}

func (i *ingress) GetIngresstDetail(ingressName, namespace string) (ingress *nwv1.Ingress, err error) {
	ingress, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress详情失败, " + err.Error()))
		return nil, errors.New("获取Ingress详情失败, " + err.Error())
	}

	return ingress, nil
}

func (i *ingress) CreateIngress(data *IngressCreate) (err error) {
	var ingressRules []nwv1.IngressRule
	var httpIngressPATHs []nwv1.HTTPIngressPath
	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Status: nwv1.IngressStatus{},
	}

	for key, value := range data.Hosts {
		ir := nwv1.IngressRule{
			Host: key,
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths: nil},
			},
		}
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: getServiceName(httpPath.ServiceName),
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			httpIngressPATHs = append(httpIngressPATHs, hip)
		}
		ir.IngressRuleValue.HTTP.Paths = httpIngressPATHs
		ingressRules = append(ingressRules, ir)
	}
	ingress.Spec.Rules = ingressRules
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Ingress失败, " + err.Error()))
		return errors.New("创建Ingress失败, " + err.Error())
	}

	return nil
}

func (i *ingress) DeleteIngress(ingressName, namespace string) (err error) {
	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Ingress失败, " + err.Error()))
		return errors.New("删除Ingress失败, " + err.Error())
	}

	return nil
}

func (i *ingress) UpdateIngress(namespace, content string) (err error) {
	var ingress = &nwv1.Ingress{}

	err = json.Unmarshal([]byte(content), ingress)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新ingress失败, " + err.Error()))
		return errors.New("更新ingress失败, " + err.Error())
	}
	return nil
}

func (i *ingress) toCells(std []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ingressCell(std[i])
	}
	return cells
}

func (i *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingresss := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		ingresss[i] = nwv1.Ingress(cells[i].(ingressCell))
	}

	return ingresss
}
