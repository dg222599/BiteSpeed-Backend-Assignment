package models

import (
	database "github.com/dg222599/BiteSpeed-Backend-Assignment/database"
)

type User struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

func (user *User) Save() (*User,error) {
 
	  err := database.DB.Create(&user).Error
	  if err!=nil{
		return &User{},err
	  }

	  return user,nil
}