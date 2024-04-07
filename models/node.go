package models

type Node struct {
	ID         int
	Link       string
	Name       string
	CreateDate string
}

// Add 添加节点
func (node *Node) Add() error {
	return DB.Create(node).Error
}

// 更新节点
func (node *Node) Update() error {
	return DB.Where("id = ?", node.ID).Updates(node).Error
}

// 查找节点是否重复
func (node *Node) Find() error {
	return DB.Where("link = ? or name = ?", node.Link, node.Name).Find(node).Error
}

// 节点列表
func (node *Node) List() ([]Node, error) {
	var nodes []Node
	err := DB.Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// 删除节点
func (node *Node) Del() error {
	return DB.Delete(node).Error
}
