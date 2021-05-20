package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uint   `gorm:"column:id;primary_key" `
	UserName          string `form:"username" gorm:"username"`
	Password          string `form:"password" gorm:"password"`
	PhoneNumber       string `form:"phone_number" gorm:"phone_number"`
	Email             string `form:"email" gorm:"email"`
	Name              string `form:"name" gorm:"name"`
	CertificateType   uint   `form:"certificate_type" gorm:"certificate_type"`
	CertificateNumber string `form:"certificate_number" gorm:"certificate_number"`
	PassengerType     uint   `form:"passenger_type" gorm:"passenger_type"`
}

func (u *User) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(passwordHash)
	return nil
}

type Passanger struct {
	ID                uint   `gorm:"column:id;primary_key" `
	UserID            uint   `form:"user_id" gorm:"user_id"`
	PhoneNumber       string `form:"phone_number" gorm:"phone_number"`
	Name              string `form:"name" gorm:"name"`
	CertificateType   uint   `form:"certificate_type" gorm:"certificate_type"`
	CertificateNumber string `form:"certificate_number" gorm:"certificate_number"`
	PassengerType     uint   `form:"passenger_type" gorm:"passenger_type"`
}
