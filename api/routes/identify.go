package routes

import (
	"github.com/dg222599/BiteSpeed-Backend-Assignment/controllers"
	"github.com/gin-gonic/gin"
)

func IdentityApp(router *gin.Engine){
	router.GET("/identify",controllers.IdentifyController)
	router.POST("/identify",controllers.LinkIdentity)
}