package users

//传给前端的数据
type UserDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user User) UserDto {
	return UserDto{
		Name: user.UserName,
	}
}
