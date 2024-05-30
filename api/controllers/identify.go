package controllers

import (
	"net/http"

	database "github.com/dg222599/BiteSpeed-Backend-Assignment/database"
	models "github.com/dg222599/BiteSpeed-Backend-Assignment/models"
	"github.com/gin-gonic/gin"
)

type User struct {
	 PhoneNumber string
	 Email string
}

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
   
	 var userDetails User

	 if err:= context.BindJSON(&userDetails) ; err!=nil{
		context.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return 
	 }
	
	 var sameEmailContacts,samePhoneContacts []models.Contact

	 matchedPhoneResults := database.DB.Where("phone_number=? and link_precedence=?",userDetails.PhoneNumber,"primary").Find(&samePhoneContacts)
	 matchedEmailResults := database.DB.Where("email=? and link_precedence=?",userDetails.Email,"primary").Find(&sameEmailContacts)

	 if matchedEmailResults.RowsAffected > 0 && matchedPhoneResults.RowsAffected > 0 {
		 // found matched records with both email and phone
	 } else if matchedEmailResults.RowsAffected > 0 || matchedPhoneResults.RowsAffected > 0 {
		 // found matched records based on either the Phone or Email
		 var parentContactID uint
		 if matchedPhoneResults.RowsAffected > 0 { 

			//need to check first if the record is new or not
			if userDetails.Email == sameEmailContacts[0].Email {
				 //Email and Phone are same so not a new record
				 context.JSON(http.StatusContinue,gin.H{"status":"record already present"})
				 return
			}	
			
			parentContactID = samePhoneContacts[0].ID
		 } else {

			if userDetails.PhoneNumber == samePhoneContacts[0].PhoneNumber {
				//Email and Phone are same so not a new record
				context.JSON(http.StatusContinue,gin.H{"status":"record already present"})
				return
		   }	
			 
			 parentContactID = sameEmailContacts[0].ID
		 }
		
		newContact := models.Contact{
			PhoneNumber: userDetails.PhoneNumber,
			Email:userDetails.Email,
			LinkedID:parentContactID,
			LinkPrecedence : "secondary",
			
		}

		saveResult := database.DB.Save(&newContact)

		if saveResult.RowsAffected > 0 {
			context.JSON(http.StatusOK,gin.H{"email":newContact.Email,"phoneNumber":newContact.PhoneNumber})
		} else {
			context.JSON(http.StatusBadRequest,gin.H{"status":"RECORD NOT SAVED" ,"error":saveResult.Error})
		}
		  

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
			context.JSON(http.StatusBadRequest,gin.H{"status":"RECORD NOT SAVED" ,"error":saveResult.Error})
		 }
	 }
}