package users

import (
	"log"
	"net/http"

	"12306.com/12306/common"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {
	var registerPara = User{}
	if err := ctx.ShouldBind(&registerPara); err != nil {
		common.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "")
		return
	}
	// 获取参数
	userName := registerPara.UserName
	password := registerPara.Password
	// 数据验证

	// 判断手机号是否存在
	DB := common.GetDB()
	if isUserNameExist(DB, userName) {
		common.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}
	// 创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		common.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := User{
		UserName: userName,
		Password: string(hasedPassword),
	}
	DB.Create(&newUser)

	// 发放token
	token, err := common.ReleaseToken(newUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	common.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = User{}
	ctx.Bind(&requestUser)
	// 获取参数
	password := requestUser.Password
	username := requestUser.UserName

	if len(password) < 6 {
		common.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 判断手机号是否存在
	var user User
	DB.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	// 判断密码收否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	common.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": ToUserDto(user.(User))}})
}

func isUserNameExist(db *gorm.DB, userName string) bool {
	var user User
	db.Where("userName = ?", userName).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}
