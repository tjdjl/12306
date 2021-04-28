package users

import (
	"github.com/gin-gonic/gin"
)

type UserModelValidator struct {
	UserPara struct {
		Username          string `form:"username" json:"username" binding:"required,alphanum,min=4,max=255"`
		Email             string `form:"email" json:"email" binding:"required,email"`
		Password          string `form:"password" json:"password" binding:"required,min=8,max=255"`
		PhoneNumber       string `form:"phoneNumber" json:"phoneNumber" binding:"required"`
		Name              string `form:"name" json:"name" binding:"required"`
		CertificateType   uint   `form:"certificateType" json:"certificateType" binding:"required"`
		CertificateNumber string `form:"certificateNumber" json:"certificateNumber" binding:"required"`
		PassengerType     uint   `form:"passengerType" gorm:"passengerType" binding:"required"`
	} `json:"user"`
	User User `json:"-"`
}

func (self *UserModelValidator) Bind(c *gin.Context) error {
	err := c.Bind(self)
	if err != nil {
		return err
	}
	self.User.UserName = self.UserPara.Username
	self.User.Email = self.UserPara.Email
	self.User.PhoneNumber = self.UserPara.PhoneNumber
	self.User.Name = self.UserPara.Name
	self.User.CertificateType = self.UserPara.CertificateType
	self.User.CertificateNumber = self.UserPara.CertificateNumber
	self.User.PassengerType = self.UserPara.PassengerType

	self.User.setPassword(self.User.Password)

	return nil
}

// You can put the default value of a Validator here
func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	return userModelValidator
}
