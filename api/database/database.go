package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnectDB() {
	DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil{
		 panic(err)
	} else {
		fmt.Println("Connected to DB Successfully!!",DB)
	}
}