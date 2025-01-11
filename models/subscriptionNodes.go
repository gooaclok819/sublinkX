package models

type SubscriptionNodes struct {
	SubscriptionID int `gorm:"primaryKey"`
	NodeID         int `gorm:"primaryKey"`
	Sort           int `gorm:"default:0"`
}

// 更新节点列表建立多对多关系并更新排序
func (sn *SubscriptionNodes) UpdateNodesWithSort(subID int, nodeSorts []SubscriptionNodes) error {
	for _, node := range nodeSorts {
		subNode := SubscriptionNodes{
			SubscriptionID: subID,
			NodeID:         node.NodeID,
			Sort:           node.Sort, // 使用传入的排序字段
		}
		// 使用 Save 方法进行插入或更新
		if err := DB.Save(&subNode).Error; err != nil {
			return err
		}
	}
	return nil
}

// 根据订阅ID删除关联的节点排序信息的方法
func (sn *SubscriptionNodes) DeleteNodesBySubscriptionID(subID int) error {
	return DB.Where("subscription_id = ?", subID).Delete(&SubscriptionNodes{}).Error
}

// 根据节点ID删除节点排序信息的方法
func (sn *SubscriptionNodes) DeleteNodeSortByNodeID(nodeID int) error {
	return DB.Where("node_id = ?", nodeID).Delete(&SubscriptionNodes{}).Error
}
