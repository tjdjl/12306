package users

//用于注册时：参数绑定   +   访问数据库
type User struct {
	ID        uint `gorm:"primary_key" `
	UserName string `form:"username" gorm:"type:varchar(20);not null"`
	Password string `form:"password" gorm:"size:255;not null"`
	//Telephone string `form:"telephone" gorm:"type:varchar(110);not null;unique"`
}







