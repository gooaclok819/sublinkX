package models

type Subcription struct {
	ID         int
	Name       string
	Config     string   `gorm:"embedded"`
	Nodes      []Node   `gorm:"many2many:subcription_nodes;"` // 多对多关系
	IPLogs     []IPLogs `gorm:"foreignKey:SubcriptionID;"`    // 一对多关系 约束父表被删除子表记录跟着删除
	CreateDate string
}

// Add 添加订阅
func (sub *Subcription) Add() error {
	return DB.Create(sub).Error
}

// 更新订阅
func (sub *Subcription) Update() error {
	return DB.Where("id = ?", sub.ID).Updates(sub).Error
}

// 查找订阅
func (sub *Subcription) Find() error {
	return DB.Where("id = ? or name = ?", sub.ID, sub.Name).Find(sub).Error
}

// 订阅列表
func (sub *Subcription) List() ([]Subcription, error) {
	var subs []Subcription
	err := DB.Find(&subs).Error
	if err != nil {
		return nil, err
	}
	for i := range subs {
		DB.Model(&subs[i]).Association("Nodes").Find(&subs[i].Nodes)
	}

	// DB.Model(sub).Association("Nodes").Find(&subs)
	return subs, nil
}

// 删除订阅
func (sub *Subcription) Del() error {
	err := DB.Model(sub).Association("Nodes").Clear()
	if err != nil {
		return err
	}
	return DB.Delete(sub).Error
}

// 添加节点列表建立多对多关系
func (sub *Subcription) AddNode() error {
	return DB.Model(sub).Association("Nodes").Append(sub.Nodes)
}
