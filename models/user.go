package models

type User struct {
	ID       int
	Username string
	Password string
	Role     string
	Nickname string
}

func (user *User) Create() error { // 创建用户
	return DB.Create(user).Error
}
func (user *User) Set(UpdateUser *User) error { // 设置用户
	return DB.Where("username = ?", user.Username).Updates(UpdateUser).Error
}
func (user *User) Verify() error { // 验证用户
	return DB.Where("username = ? AND password = ?", user.Username, user.Password).First(user).Error
}

func (user *User) Find() error { // 查找用户
	return DB.Where("username = ? ", user.Username).First(user).Error
}

func (user *User) All() ([]User, error) { // 获取所有用户
	var users []User
	err := DB.Find(&users).Error
	return users, err
}

func (user *User) Del() error { // 删除用户
	return DB.Delete(user).Error
}
