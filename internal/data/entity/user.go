package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string         `gorm:"type:varchar(100);not null" json:"name"`
	Email         string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password      string         `gorm:"type:varchar(255);not null" json:"-"`
	Role          string         `gorm:"type:varchar(20);not null;default:'staff'" json:"role"`
	OTP           string         `gorm:"type:varchar(6)" json:"-"`
	OTPExpiration *gorm.DeletedAt `gorm:"index" json:"-"`
}
