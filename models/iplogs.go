package models

type SubLogs struct {
	ID             int
	IP             string
	Date           string
	Addr           string
	Count          int
	SubscriptionID int
}

// Add 添加IP
func (ipLog *SubLogs) Add() error {
	return DB.Create(ipLog).Error
}

// Find 查找IP
func (ipLog *SubLogs) Find(id int) error {
	return DB.Where("ip = ? or subscription_id  = ?", ipLog.IP, id).First(ipLog).Error
}

// Update 更新IP
func (ipLog *SubLogs) Update() error {
	return DB.Where("id = ? or ip = ?", ipLog.ID, ipLog.IP).Updates(ipLog).Error
}

// List 获取IP列表
func (ipLog *SubLogs) List() ([]SubLogs, error) {
	var iplogs []SubLogs
	err := DB.Find(&iplogs).Error
	if err != nil {
		return nil, err
	}
	return iplogs, nil
}

// Del 删除IP
func (ipLog *SubLogs) Del() error {
	return DB.Delete(ipLog).Error
}
