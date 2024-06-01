package models

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	PhoneNumber *string `json:"phoneNumber" gorm:"column:phone_number"`
	Email *string `json:"email"`
	LinkedID uint `json:"linkedID" gorm:"column:linked_id"`
	LinkPrecedence string `json:"linkPrecedence" gorm:"column:link_precedence"`
	
}