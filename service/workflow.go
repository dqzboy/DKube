package service

import (
	"dkube/dao"
	"dkube/model"
)

var Workflow workflow

type workflow struct{}

type WorkflowCreate struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Replicas      int32                  `json:"replicas"`
	Image         string                 `json:"image"`
	Label         map[string]string      `json:"label"`
	Cpu           string                 `json:"cpu"`
	Memory        string                 `json:"memory"`
	ContainerPort int32                  `json:"container_port"`
	HealthCheck   bool                   `json:"health_check"`
	HealthPath    string                 `json:"health_path"`
	Type          string                 `json:"type"`
	Port          int32                  `json:"port"`
	NodePort      int32                  `json:"node_port"`
	Hosts         map[string][]*HttpPath `json:"hosts"`
}

func (w *workflow) GetList(name, namespace string, page, limit int) (data *dao.WorkflowResp, err error) {
	data, err = dao.Workflow.GetWorkflows(name, namespace, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *workflow) GetById(id int) (data *model.Workflow, err error) {
	data, err = dao.Workflow.GetById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *workflow) CreateWorkFlow(data *WorkflowCreate) (err error) {
	var ingressName string
	if data.Type == "Ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}

	workflow := &model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    getServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}
	err = dao.Workflow.Add(workflow)
	if err != nil {
		return err
	}
	err = createWorkflowRes(data)
	if err != nil {
		return err
	}
	return err
}

func (w *workflow) DelById(id int) (err error) {
	workflow, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}
	err = delWorkflowRes(workflow)
	if err != nil {
		return err
	}
	err = dao.Workflow.DelById(id)
	if err != nil {
		return err
	}
	return
}

func delWorkflowRes(workflow *model.Workflow) (err error) {
	err = Deployment.DeleteDeployment(workflow.Name, workflow.Namespace)
	if err != nil {
		return err
	}
	err = Servicev1.DeleteService(getServiceName(workflow.Name), workflow.Namespace)
	if err != nil {
		return err
	}

	if workflow.Type == "Ingress" {
		err = Ingress.DeleteIngress(getIngressName(workflow.Name), workflow.Namespace)
		if err != nil {
			return err
		}
	}

	return nil
}

func createWorkflowRes(data *WorkflowCreate) (err error) {
	dc := &DeployCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Replicas:      data.Replicas,
		Image:         data.Image,
		Label:         data.Label,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
		ContainerPort: data.ContainerPort,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
	}
	err = Deployment.CreateDeployment(dc)
	if err != nil {
		return err
	}
	var serviceType string
	if data.Type != "Ingress" {
		serviceType = data.Type
	} else {
		serviceType = "ClusterIP"
	}
	sc := &ServiceCreate{
		Name:          getServiceName(data.Name),
		Namespace:     data.Namespace,
		Type:          serviceType,
		ContainerPort: data.ContainerPort,
		Port:          data.Port,
		NodePort:      data.NodePort,
		Label:         data.Label,
	}
	if err := Servicev1.CreateService(sc); err != nil {
		return err
	}
	var ic *IngressCreate
	if data.Type == "Ingress" {
		ic = &IngressCreate{
			Name:      getIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
		err = Ingress.CreateIngress(ic)
		if err != nil {
			return err
		}
	}

	return nil
}

func getServiceName(workflowName string) (serviceName string) {
	return workflowName + "-svc"
}

func getIngressName(workflowName string) (ingressName string) {
	return workflowName + "-ing"
}
