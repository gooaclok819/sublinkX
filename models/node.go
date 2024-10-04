package models

import (
	"errors"
	"gorm.io/gorm"
)

type Node struct {
	ID         int
	Link       string
	Name       string
	Sort       int `gorm:"default:0"`
	CreateDate string
}

// Add 添加节点
func (node *Node) Add() error {
	var maxSort int
	DB.Model(&Node{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)
	node.Sort = maxSort + 1

	return DB.Create(node).Error
}

// Update 更新节点
func (node *Node) Update() error {
	return DB.Model(node).Updates(node).Error
}

// UpdateById 根据 ID 更新节点
func (node *Node) UpdateById() error {
	// 检查是否有其他节点具有相同的名称或链接
	existingNode, err := node.FindByLinkOrName()
	if err != nil {
		return err
	}
	if existingNode != nil {
		return errors.New("节点名称或链接已存在")
	}
	return DB.Model(node).Where("id = ?", node.ID).Updates(node).Error
}

// UpdateSort 批量更新节点的排序
func UpdateSort(nodes []Node) error {
	tx := DB.Begin()
	for _, node := range nodes {
		err := tx.Model(&Node{}).Where("id = ?", node.ID).Update("sort", node.Sort).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// Find 根据链接或者名称查找节点
func (node *Node) Find() error {
	return DB.Where("link = ? or name = ?", node.Link, node.Name).First(node).Error
}

// FindByLinkOrName 根据名称或链接查找节点
func (node *Node) FindByLinkOrName() (*Node, error) {
	var existingNode Node
	err := DB.Where("name = ? OR link = ?", node.ID, node.Name, node.Link).First(&existingNode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有找到记录，返回 nil 和 gorm.ErrRecordNotFound
			return nil, nil
		}
		// 其他错误直接返回
		return nil, err
	}
	return &existingNode, nil
}

// FindByID 根据 ID 查找节点
func (node *Node) FindByID(id int) error {
	return DB.Where("id = ?", id).First(node).Error
}

// List 节点列表
func (node *Node) List() ([]Node, error) {
	var nodes []Node
	err := DB.Order("sort ASC").Order("create_date DESC").Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// Del 删除节点
func (node *Node) Del() error {
	return DB.Delete(node).Error
}
