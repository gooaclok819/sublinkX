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
	return DB.Model(node).Updates(node).Error
}

// 查找节点是否重复
func (node *Node) Find() error {
	return DB.Where("link = ? or name = ?", node.Link, node.Name).First(node).Error
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
	// 删除节点排序信息
	sn := &SubscriptionNodes{}
	err := sn.DeleteNodeSortByNodeID(node.ID)
	if err != nil {
		return err
	}

	// 删除节点
	return DB.Delete(node).Error
}
