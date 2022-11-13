package user

type User struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Email    string
	Birthday string
	Age      int
}

func (User) TableName() string {
	return "t_user"
}

// 注册请求数据结构
type SigninRequest struct {
	Username string
	Email    string
}

// 向 user 表中插入一行数据（包括用户名、邮件地址等基本信息）
func (u *User) signin() error {
	db.Select("Email", "Username").Create(u)
	return nil
}
