package users

import (
	"log"
	"net/http"

	"12306.com/12306/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	//校验数据
	// userModelValidator := NewUserModelValidator()
	// if err := userModelValidator.Bind(c); err != nil {
	// 	c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 0, "msg": "参数错误", "data": err.Error()})
	// 	return
	// }
	// user := userModelValidator.User
	//获取参数
	var user = User{}
	c.Bind(&user)
	user.setPassword(user.Password)
	//在数据库中创建用户
	DB := common.GetDB()
	err := DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 0, "msg": "用户创建失败", "data": err.Error()})
		return
	}
	//生成token
	token, err := common.ReleaseToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 0, "msg": "系统异常", "data": ""})
		log.Printf("token generate error : %v", err)
		return
	}
	//返回结果
	data := gin.H{"token": token}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "成功", "data": data}) //返回结果
}

func Login(c *gin.Context) {
	// 获取参数
	var requestUser = User{}
	c.Bind(&requestUser)
	password := requestUser.Password
	username := requestUser.UserName
	// 按照用户名从数据库取得用户
	DB := common.GetDB()
	var user User
	DB.Where("user_name = ?", username).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}
	// 生成token
	token, err := common.ReleaseToken(user.ID)
	if err != nil {
		log.Printf("token generate error : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		return
	}
	// 返回结果
	data := gin.H{"token": token}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "成功", "data": data})
}

func AddPassenger(c *gin.Context) {
	var passanger = Passanger{}
	c.Bind(&passanger)
	user, _ := c.Get("user")
	passanger.UserID = user.(User).ID
	// 向数据库中增加乘客
	DB := common.GetDB()
	err := DB.Create(&passanger).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "成功"})
}

func UpdatePassenger(c *gin.Context) {
}

func QueryPassenger(c *gin.Context) {
	// 获取用户信息
	user, _ := c.Get("user")

	// 按照userID从数据库取得乘车人
	DB := common.GetDB()
	var passangers []Passanger
	err := DB.Where("user_id = ?", user.(User).ID).Find(&passangers).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 422, "msg": err.Error()})
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "成功", "data": passangers})
}
