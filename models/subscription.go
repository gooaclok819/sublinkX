package models

type Subscription struct {
	ID         int
	Name       string
	Config     string    `gorm:"embedded"`
	Nodes      []Node    `gorm:"many2many:subscription_nodes;"` // 多对多关系
	SubLogs    []SubLogs `gorm:"foreignKey:SubscriptionID;"`    // 一对多关系 约束父表被删除子表记录跟着删除
	CreateDate string
}

// Add 添加订阅
func (sub *Subscription) Add() error {
	return DB.Create(sub).Error
}

// AddNode 添加节点列表建立多对多关系
func (sub *Subscription) AddNode() error {
	return DB.Model(sub).Association("Nodes").Append(sub.Nodes)
}

// Update 更新订阅
func (sub *Subscription) Update() error {
	return DB.Where("id = ? or name = ?", sub.ID, sub.Name).Updates(sub).Error
}

// UpdateNodes 更新节点列表建立多对多关系
func (sub *Subscription) UpdateNodes() error {
	return DB.Model(sub).Association("Nodes").Replace(sub.Nodes)
}

// Find 查找订阅
func (sub *Subscription) Find() error {
	return DB.Where("id = ? or name = ?", sub.ID, sub.Name).First(sub).Error
}

// GetSub 读取订阅
func (sub *Subscription) GetSub() error {
	// err := DB.Find(sub).Error
	// if err != nil {
	// 	return err
	// }
	return DB.Model(sub).Association("Nodes").Find(&sub.Nodes)
}

// List 订阅列表
func (sub *Subscription) List() ([]Subscription, error) {
	var subs []Subscription
	err := DB.Find(&subs).Error
	if err != nil {
		return nil, err
	}
	for i := range subs {
		err := DB.Model(&subs[i]).Association("Nodes").Find(&subs[i].Nodes)
		if err != nil {
			return nil, err
		}
		logsErr := DB.Model(&subs[i]).Association("SubLogs").Find(&subs[i].SubLogs)
		if logsErr != nil {
			return nil, logsErr
		}
	}
	return subs, nil
}

// IPLogUpdate IP记录更新
func (sub *Subscription) IPLogUpdate() error {
	return DB.Model(sub).Association("SubLogs").Replace(&sub.SubLogs)
}

// Del 删除订阅
func (sub *Subscription) Del() error {
	err := DB.Model(sub).Association("Nodes").Clear()
	if err != nil {
		return err
	}
	return DB.Delete(sub).Error
}
