package models

type IPLogs struct {
	ID            int
	IP            string
	Date          string
	SubcriptionID int
}

// Add 添加IP
func (iplog *IPLogs) Add() error {
	return DB.Create(iplog).Error
}

// Update 更新IP
func (iplog *IPLogs) Update() error {
	return DB.Where("id = ?", iplog.ID).Updates(iplog).Error
}

// List 获取IP列表
func (iplog *IPLogs) List() ([]IPLogs, error) {
	var iplogs []IPLogs
	err := DB.Find(&iplogs).Error
	if err != nil {
		return nil, err
	}
	return iplogs, nil
}

// Del 删除IP
func (iplog *IPLogs) Del() error {
	return DB.Delete(iplog).Error
}
