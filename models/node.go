package models

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type GroupNode struct {
	gorm.Model
	ID    int
	Name  string
	Nodes []Node `gorm:"many2many:group_node_nodes"` // 多对多关联字段
}

type Node struct {
	gorm.Model
	ID         int
	Name       string
	Link       string
	GroupNodes []GroupNode `gorm:"many2many:group_node_nodes"` // 反向关联字段
}

// 创建分组
func (gn *GroupNode) Add() error {
	// 检查分组是否已存在
	var existingGroup GroupNode
	result := DB.Model(gn).Where("name = ?", gn.Name).First(&existingGroup) // 查询数据库中是否存在同名的分组
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(result.Error)
		return result.Error // 如果查询出错，返回错误
	}
	if result.RowsAffected > 0 { // 如果查询到分组已存在
		log.Println("分组已存在")
		return nil // 不返回错误存在就跳过

	}
	return DB.FirstOrCreate(gn, GroupNode{Name: gn.Name}).Error
}

// 关联分组
func (gn *GroupNode) Ass(n *Node) error {
	result := DB.Model(gn).Where("name = ?", gn.Name).First(gn) // 查找分组
	log.Println(gn)
	if result.Error != nil {
		log.Println(result.Error)
	}
	result = DB.Model(n).Where("name = ?", n.Name).First(n) // 查找节点
	if result.Error != nil {
		log.Println(result.Error)
	}
	return DB.Model(gn).Association("Nodes").Append(n)
}

// 更新分组信息
func (gn *GroupNode) Update(NewGn *GroupNode) error {
	// 读取分组数据
	var FirstGn GroupNode
	result := DB.Model(gn).Where("id = ? or name = ?", NewGn.ID, NewGn.Name).First(&FirstGn)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	if result.RowsAffected > 0 {
		return errors.New("分组已存在")
	}
	return DB.Model(gn).Where("id = ? or name = ?", gn.ID, gn.Name).Updates(&NewGn).Error
}

// 删除分组
func (gn *GroupNode) Del() error {
	// 读取分组数据
	result := DB.Model(gn).Where("id = ? or name = ?", gn.ID, gn.Name).First(&gn)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	// 读取分组关联的节点数据
	result = DB.Model(gn).Preload("Nodes").First(gn) // 预加载分组关联的节点数据
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	log.Println(gn.Nodes)
	err := DB.Model(gn).Association("Nodes").Delete(gn.Nodes)
	if err != nil {
		log.Println("解除关联失败", err)
		return err
	}
	return DB.Model(gn).Where("id = ? or name = ?", gn.ID, gn.Name).Delete(gn).Error // 删除分组记录
}

// 查看所有分组

func GetGroupNodeList() ([]GroupNode, error) {
	var gns []GroupNode
	result := DB.Model(gns).Preload("Nodes").Find(&gns)
	if result.Error != nil {
		return nil, errors.New("没有任何分组")
	}
	return gns, result.Error
}

/* 下面为节点的增删改查 */

// 添加节点的方法
func (n *Node) Add() error {
	// 检查节点是否已存在
	var existingNode Node
	result := DB.Model(n).Where("link = ? and name =?", n.Link, n.Name).First(&existingNode) // 查询数据库中是否存在同名同链接的节点
	if result.RowsAffected > 0 {
		log.Println(result.RowsAffected)
		return errors.New("节点已经存在") // 如果查询出错，返回错误
	}
	return DB.Model(n).Create(n).Error // 使用 GORM 创建新的节点记录
}

// 删除节点
func (n *Node) Del() error {
	// 查看是否有关联
	if len(n.GroupNodes) > 0 {
		err := DB.Model(n).Association("GroupNodes").Delete(n.GroupNodes)
		if err != nil {
			return err
		}
	}

	// Unscoped  硬删除
	// 默认删除是软删除 数据库仍然存在记录
	return DB.Model(n).Delete(n).Error
}

// 更新节点

func (n *Node) UpdateNode(New *Node) error {
	// 检查节点是否已存在
	var n1 Node
	result := DB.Model(n).Where("link = ? or name =?", New.Link, New.Name).First(&n1) // 查询数据库中是否存在同名同链接的节点
	if result.RowsAffected > 0 {
		log.Println("节点已经存在", result.Error)
		return errors.New("节点已经存在") // 如果查询出错，返回错误
	}
	// 更新记录
	return DB.Model(n).Where("id = ?", n.ID).Updates(New).Error
}

// 更新关联分组

func (n *Node) UpdateGroup(NewGroups []GroupNode) error {
	// 检测节点是否存在
	result := DB.Model(n).Where("id = ? or name = ?", n.ID, n.Name).First(&n) // 查找节点
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error // 如果查询出错，返回错误
	}

	// 检查分组是否已存在
	var NewGroupDatas []GroupNode

	for _, NewGroup := range NewGroups {

		var NewGroupData GroupNode

		if NewGroup.Name == "" {
			log.Println("分组名称不能为空", NewGroup)
			return errors.New("分组名称不能为空") // 如果分组名称为空，返回错误
		}
		result := DB.Model(GroupNode{}).Where("name = ?", NewGroup.Name).First(&NewGroupData) // 查找分组并将数据放在NewGroupData
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error // 如果查询出错，返回错误
		}
		NewGroupDatas = append(NewGroupDatas, NewGroupData) // 将新的数据添加到更新列表中
	}

	// 更新记录
	return DB.Model(n).Association("GroupNodes").Replace(NewGroupDatas) // 替换分组节点
}

// 查看所有节点

func GetNodeList() ([]Node, error) {
	var ns []Node
	result := DB.Model(ns).Preload("GroupNodes").Find(&ns)
	if result.Error != nil {
		return nil, result.Error
	}
	return ns, result.Error
}
