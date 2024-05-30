package main

import (
	database "github.com/dg222599/BiteSpeed-Backend-Assignment/database"
	models "github.com/dg222599/BiteSpeed-Backend-Assignment/models"
	routes "github.com/dg222599/BiteSpeed-Backend-Assignment/routes"
	"github.com/gin-gonic/gin"
)

func main(){ 
	
	database.ConnectDB()

	
	if err:=database.DB.AutoMigrate(&models.Contact{}) ; err!=nil{
		
		panic(err)
	}
	// if err:= database.DB.AutoMigrate(&models.User{}) ; err!=nil{
	// 	 panic(err)
	// }
    
	router := gin.Default()
	
	
	routes.IdentityApp(router)

	router.Run(":3000")
}