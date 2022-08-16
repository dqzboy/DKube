package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct{}

//定义列表的返回内容，Items是deployment元素列表，Total为deployment元素数量
type DeploymentsResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int
}

//定义DeployCreate结构体，用于创建deployment需要的参数属性的定义
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
}

//定义DeploysNp类型,用于返回namespace中的deployment数据
type DeploysNp struct {
	Namespace string `json:"namespace"`
	DeployNum int    `json:"deployment_num"`
}

//获取deployment列表，支持过滤、排序、分页
func (d *deployment) GetDeployments(filterName, namespace string, limit, page int) (deploymentsResp *DeploymentsResp, err error) {
	//获取deploymentList类型的deployment列表
	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//打印日志,排错
		logger.Error("获取Deployment列表失败," + err.Error())
		//返回给上一层，最终返回给前端，前端打印出该error
		return nil, errors.New("获取Deployment列表失败," + err.Error())
	}

	//将deploymentList中的deployment列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
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

	//将DataCell类型的deployment列表转换成appsv1.deployment列表
	deployments := d.fromCells(data.GenericDataList)

	return &DeploymentsResp{
		Items: deployments,
		Total: total,
	}, nil

}

//获取Deployment详情
func (d *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Deployment详情失败," + err.Error())
		return nil, errors.New("获取Deployment详情失败," + err.Error())
	}
	return deployment, nil
}

//修改Deployment副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error) {
	//获取autoscalingv1.Scale类型的对象，能点出当前的副本数
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment副本数信息失败," + err.Error()))
		return 0, errors.New("获取Deployment副本数信息失败," + err.Error())
	}
	//修改副本数
	scale.Spec.Replicas = int32(scaleNum)
	//更新副本数，传入scale对象
	newScale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment副本数信息失败," + err.Error()))
		return 0, errors.New("更新Deployment副本数信息失败," + err.Error())
	}
	return newScale.Spec.Replicas, nil
}

//创建Deployment,接收DeployCreate对象
func (d *deployment) CreateDeployment(data *DeployCreate) (err error) {
	//初始化appsv1.Deployment类型的对象，并将入参的data数据放进去
	deployment := &appsv1.Deployment{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		//Spec中定义副本数、选择器、以及Pod属性
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				//定义Pod名和标签
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				//定义容器名、镜像和端口
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: data.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		//Status定义资源的运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status: appsv1.DeploymentStatus{},
	}
	//判断是否打开健康检查功能，如果打开，则定义ReadinessProbe和LivenessProbe
	if data.HealthCheck {
		//设置第一个容器的ReadnessProbe，因为我们pod中只有一个容器，所以直接使用index 0 即可
		//如果pod中有多个容器，则这里需要使用for循环去定义
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
					//Type=0则表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
					//Type=1则表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			//初始化等待时间
			InitialDelaySeconds: 5,
			//超时时间
			TimeoutSeconds: 5,
			//执行间隔
			PeriodSeconds: 5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			//初始化等待时间
			InitialDelaySeconds: 15,
			//超时时间
			TimeoutSeconds: 5,
			//执行间隔
			PeriodSeconds: 5,
		}
	}
	//定义容器的limit和request资源
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	//调用SDK去更新Deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(data.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Deployment失败," + err.Error()))
		return errors.New("创建Deployment失败," + err.Error())
	}
	return nil
}

//删除deployment
func (d *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Deployment失败," + err.Error())
		return errors.New("删除Deployment失败," + err.Error())
	}
	return nil
}

//重启deployment
func (d *deployment) RestartDeployment(deploymentName, namespace string) (err error) {
	//此功能等同于kubectl命令
	//使用patchData Map组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	//序列化为字节，因为patch方法只接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(errors.New("josn序列化失败," + err.Error()))
		return errors.New("json序列化失败," + err.Error())
	}
	//调用patch方法更新deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启Deployment失败," + err.Error()))
		return errors.New("重启Deployment失败," + err.Error())
	}
	return
}

//更新deployment
func (d *deployment) UpdateDeployment(namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error("反序列化失败," + err.Error())
		return errors.New("反序列化失败," + err.Error())
	}
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Deployment失败," + err.Error())
		return errors.New("更新Deployment失败," + err.Error())
	}
	return nil
}

//获取每个namespace的deployment数量
func (d *deployment) GetDeployNumPerNp() (deploysNps []*DeploysNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		deploysNps := &DeploysNp{
			Namespace: namespace.Name,
			DeployNum: len(deploymentList.Items),
		}
		deploysNps = append(deploysNps, deploysNp)
	}
	return deploysNps, nil
}

//类型转换的方法
func (d *deployment) toCells(deployments []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(deployments))
	for i := range deployments {
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}

func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		//cells[i].(podCell)是将DataCell类型转成podCell类型
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}
