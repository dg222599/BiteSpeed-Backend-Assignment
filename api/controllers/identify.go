package controllers

import (
	"net/http"

	database "github.com/dg222599/BiteSpeed-Backend-Assignment/database"
	models "github.com/dg222599/BiteSpeed-Backend-Assignment/models"
	"github.com/gin-gonic/gin"
)

func IdentifyController(c *gin.Context){
	var allContacts []models.Contact
	result := database.DB.Find(&allContacts)
	if result.Error != nil {
		 c.Error(result.Error)
	} else {
		c.JSON(http.StatusFound,gin.H{"allContacts":allContacts})
	}
	
}

func LinkIdentity(context *gin.Context){
   
	 var userDetails models.User

	 if err:= context.BindJSON(&userDetails) ; err!=nil{
		context.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return 
	 }
	
	 var sameEmailContacts,samePhoneContacts []models.Contact

	 result1 := database.DB.Where("phone_number=?",userDetails.PhoneNumber).Find(&samePhoneContacts)
	 result2 := database.DB.Where("email=?",userDetails.Email).Find(&sameEmailContacts)

	 if result1.RowsAffected > 0 {

		context.JSON(http.StatusOK,gin.H{"email":parentContact1.Email,"phoneNumber":parentContact1.PhoneNumber})
	 } else if result2.RowsAffected > 0 {

		context.JSON(http.StatusOK,gin.H{"email":parentContact2.Email,"phoneNumber":parentContact2.PhoneNumber})
	 } else {
		
		 //need to create the new contact since  there is no contact with this phone/email
		 newContact := models.Contact{
				PhoneNumber: userDetails.PhoneNumber,
				Email:userDetails.Email,
				LinkPrecedence : "primary",
		 }

		 saveResult := database.DB.Save(&newContact)
		 if saveResult.RowsAffected > 0 {
			context.JSON(http.StatusOK,gin.H{"email":newContact.Email,"phoneNumber":newContact.PhoneNumber})
		 } else {
			context.JSON(http.StatusNotFound,gin.H{"status":"RECORD NOT SAVED"})
		 }
	 }
}