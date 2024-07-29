package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	MavisModel
	Name           string `json:"name"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password" gorm:"-"`
	HashedPassword string `json:"-" gorm:"column:password"`
	Phone          string `json:"phone"`
	StoreID        *uint  `json:"store_id"`
	Store          *Store
	Status         *bool  `json:"status"`
	Token          string `json:"-"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if len(u.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = ""
		u.HashedPassword = string(hashedPassword)
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	var checkUser User
	if err := tx.Where("username = ?", u.Username).First(&checkUser).Error; err != nil {
		return nil
	} else {
		return errors.New("username already exist")
	}
}
