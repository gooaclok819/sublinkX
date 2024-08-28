package models

type SubscriptionNodes struct {
	SubscriptionID int `gorm:"primaryKey"`
	NodeID         int `gorm:"primaryKey"`
	Sort           int `gorm:"default:0"`
}
