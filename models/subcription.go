package models

import (
	"sort"
)

type Subcription struct {
	ID         int
	Name       string
	Config     string    `gorm:"embedded"`
	Nodes      []Node    `gorm:"many2many:subcription_nodes;"` // 多对多关系
	SubLogs    []SubLogs `gorm:"foreignKey:SubcriptionID;"`    // 一对多关系 约束父表被删除子表记录跟着删除
	CreateDate string
}

// Add 添加订阅
func (sub *Subcription) Add() error {
	return DB.Create(sub).Error
}

// 添加节点列表建立多对多关系
func (sub *Subcription) AddNode() error {
	return DB.Model(sub).Association("Nodes").Append(sub.Nodes)
}

// 更新订阅
func (sub *Subcription) Update() error {
	return DB.Where("id = ? or name = ?", sub.ID, sub.Name).Updates(sub).Error
}

// 更新节点列表建立多对多关系
func (sub *Subcription) UpdateNodes() error {
	tx := DB.Begin()

	// 清除旧的节点排序信息
	err := tx.Model(&SubscriptionNodes{}).Where("subscription_id = ?", sub.ID).Delete(&SubscriptionNodes{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新节点列表并建立多对多关系
	err = tx.Model(sub).Association("Nodes").Replace(sub.Nodes)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 保存新的节点排序信息
	for i, node := range sub.Nodes {
		subscriptionNode := SubscriptionNodes{
			SubscriptionID: sub.ID,
			NodeID:         node.ID,
			Sort:           i + 1,
		}
		err = tx.Create(&subscriptionNode).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// 查找订阅
func (sub *Subcription) Find() error {
	if err := DB.Where("id = ? or name = ?", sub.ID, sub.Name).First(sub).Error; err != nil {
		return err
	}
	return sub.loadNodesSorted()
}

// 读取订阅
func (sub *Subcription) GetSub() error {
	return sub.loadNodesSorted()
}

// 订阅列表
func (sub *Subcription) List() ([]Subcription, error) {
	var subs []Subcription
	err := DB.Find(&subs).Error
	if err != nil {
		return nil, err
	}
	for i := range subs {
		if err := subs[i].loadNodesSorted(); err != nil {
			return nil, err
		}
		if logsErr := DB.Model(&subs[i]).Association("SubLogs").Find(&subs[i].SubLogs); logsErr != nil {
			return nil, logsErr
		}
	}
	return subs, nil
}

func (sub *Subcription) IPlogUpdate() error {
	return DB.Model(sub).Association("SubLogs").Replace(&sub.SubLogs)
}

// 删除订阅
func (sub *Subcription) Del() error {
	// 清空关联订阅的节点排序信息
	sn := &SubscriptionNodes{}
	err := sn.DeleteNodesBySubscriptionID(sub.ID)
	if err != nil {
		return err
	}

	// 删除订阅
	return DB.Delete(sub).Error
}

// 加载排序后的节点列表
func (sub *Subcription) loadNodesSorted() error {
	var subscriptionNodes []SubscriptionNodes
	if err := DB.Where("subscription_id = ?", sub.ID).Find(&subscriptionNodes).Error; err != nil {
		return err
	}

	nodeIDToSort := make(map[int]int)
	for _, sn := range subscriptionNodes {
		nodeIDToSort[sn.NodeID] = sn.Sort
	}

	if err := DB.Model(sub).Association("Nodes").Find(&sub.Nodes); err != nil {
		return err
	}

	sort.Slice(sub.Nodes, func(i, j int) bool {
		sortI, okI := nodeIDToSort[sub.Nodes[i].ID]
		sortJ, okJ := nodeIDToSort[sub.Nodes[j].ID]
		if !okI {
			sortI = int(^uint(0) >> 1)
		}
		if !okJ {
			sortJ = int(^uint(0) >> 1)
		}
		return sortI < sortJ
	})

	return nil
}
