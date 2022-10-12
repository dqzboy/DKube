package dao

import (
	"dkube/db"
	"dkube/model"
	"errors"
	"github.com/wonderivan/logger"
)

var Workflow workflow

type workflow struct{}

type WorkflowResp struct {
	Items []*model.Workflow `json:"items"`
	Total int               `json:"total"`
}

func (w *workflow) GetWorkflows(filterName, namespace string, limit, page int) (data *WorkflowResp, err error) {
	startSet := (page - 1) * limit
	var (
		workflowList []*model.Workflow
		total        int
	)

	tx := db.GORM.
		Model(&model.Workflow{}).
		Where("name like ?", "%"+filterName+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&workflowList)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Workflow列表失败, " + tx.Error.Error())
		return nil, errors.New("获取Workflow列表失败, " + tx.Error.Error())
	}
	return &WorkflowResp{
		Items: workflowList,
		Total: total,
	}, nil
}

func (w *workflow) GetById(id int) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).First(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Workflow详情失败, " + tx.Error.Error())
		return nil, errors.New("获取Workflow详情失败, " + tx.Error.Error())
	}
	return workflow, nil
}

func (w *workflow) Add(workflow *model.Workflow) (err error) {
	tx := db.GORM.Create(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("创建Workflow失败, " + tx.Error.Error())
		return errors.New("创建Workflow失败, " + tx.Error.Error())
	}
	return nil
}

func (w *workflow) DelById(id int) (err error) {
	tx := db.GORM.Where("id = ?", id).Delete(&model.Workflow{})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Workflow详情失败, " + tx.Error.Error())
		return errors.New("获取Workflow详情失败, " + tx.Error.Error())
	}
	return nil
}
