package controllers

import (
	"net/http"
	"strconv"

	database "github.com/dg222599/BiteSpeed-Backend-Assignment/database"
	models "github.com/dg222599/BiteSpeed-Backend-Assignment/models"
	"github.com/gin-gonic/gin"
)

type User struct {
	 PhoneNumber string
	 Email string
}

type combinedContact struct {
	PrimaryContactId uint `json:"primaryContactId"`
	Emails []string	`json:"emails"`
	PhoneNumbers []string `json:"phoneNumbers"`
	SecondaryContactIds []uint `json:"secondaryContactIds"`

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

func DeleteContact(c *gin.Context){
	contactID,err:=strconv.ParseUint(c.Param("id"),10,64)
	if err!=nil{
		 c.JSON(http.StatusBadRequest,gin.H{"message":"Can not parse the ID to int"})
		 return 
	}

	var deletedContact models.Contact
	deleteStatus := database.DB.Where("id = ?",contactID).Delete(&deletedContact)

	if deleteStatus.RowsAffected > 0 {
		 c.JSON(http.StatusOK,gin.H{"status":"record deleted","contact":deletedContact})
	} else {
		 c.JSON(http.StatusBadRequest,gin.H{"status":"no such contact with given ID"})
	}
}

// helper function
func combineContacts(primaryID uint) combinedContact {

	// need to find all contacts which have linkedID = primaryID
	//   {
	// 	"contact":{
	// 		"primaryContatctId": number,
	// 		"emails": string[], // first element being email of primary contact 
	// 		"phoneNumbers": string[], // first element being phoneNumber of primary contact
	// 		"secondaryContactIds": number[] // Array of all Contact IDs that are "secondary" to the primary contact
	// 	}
	// }

	completeContactData := combinedContact{}

	//assign ID
	completeContactData.PrimaryContactId = primaryID

	var primaryContact models.Contact
	primaryContactResult := database.DB.Where("id=?",primaryID).Find(&primaryContact)
	if primaryContactResult.RowsAffected <=0 {
		 return completeContactData
	}

	//assign first email
	completeContactData.Emails = append(completeContactData.Emails,primaryContact.Email)
	completeContactData.PhoneNumbers = append(completeContactData.PhoneNumbers, primaryContact.PhoneNumber)

	// find all the contacts which are related to this primary contact
	var relatedContacts []models.Contact
	result:= database.DB.Where("linked_id=?",primaryID).Find(&relatedContacts)

	if result.RowsAffected <=0 {
		 return completeContactData
	}


	for index:=0;index<len(relatedContacts);index++ {
		 
		// need to append email,phonenumber,id in seconddarycontactIDs
		completeContactData.Emails = append(completeContactData.Emails,relatedContacts[index].Email)
		completeContactData.PhoneNumbers = append(completeContactData.PhoneNumbers,relatedContacts[index].PhoneNumber)
		completeContactData.SecondaryContactIds = append(completeContactData.SecondaryContactIds, relatedContacts[index].ID)
	}

	return completeContactData


	 
}
 
func LinkIdentity(context *gin.Context){
   
	 var userDetails User

	 if err:= context.BindJSON(&userDetails) ; err!=nil{
		context.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return 
	 }

	 var alreadyPresentRecord models.Contact

	 contactAlreadyPresent := database.DB.Where("phone_number=? and email=?",userDetails.PhoneNumber,userDetails.Email).First(&alreadyPresentRecord)

	 if contactAlreadyPresent.RowsAffected > 0 {
		context.JSON(http.StatusFound,gin.H{"status":"record already present"})
		return
	 } else {
		
	    var primaryEmailContact,primaryPhoneContact models.Contact
		
		// denotes the record where a primary contact has same email/phone as current contact
		phoneResultPrimary := database.DB.Where("phone_number=? and link_precedence=?",userDetails.PhoneNumber,"primary").Find(&primaryPhoneContact)
	 	emailResultPrimary := database.DB.Where("email=? and link_precedence=?",userDetails.Email,"primary").Find(&primaryEmailContact)
		
		if emailResultPrimary.RowsAffected > 0 && phoneResultPrimary.RowsAffected > 0 {
			  //found primary match in both(can be same contact can be not)
			 // email coming from one contact and phone coming from another contact
			 //create link either way
			 primaryEmailContact.LinkedID = primaryPhoneContact.ID
			 primaryEmailContact.LinkPrecedence = "secondary"
			 updatedResult:=database.DB.Save(&primaryEmailContact)

			 if updatedResult.RowsAffected > 0 {
				
				consolidatedContact:=combineContacts(primaryPhoneContact.ID)
				context.JSON(http.StatusOK,gin.H{"contact":consolidatedContact})
	
			 }
		} else if emailResultPrimary.RowsAffected > 0 || phoneResultPrimary.RowsAffected > 0 {

			 //found primary match in only one
			 newContact := models.Contact{
				PhoneNumber: userDetails.PhoneNumber,
				Email:userDetails.Email,
				LinkPrecedence : "secondary",
				
			 }

			 
			 if emailResultPrimary.RowsAffected > 0 {
				 newContact.LinkedID = primaryEmailContact.ID
			 } else {
				 newContact.LinkedID = primaryPhoneContact.ID
			 }

			 saveResult := database.DB.Save(&newContact)
			 if saveResult.RowsAffected <=0 {
					return
			 }

			 
			 consolidatedContact:=combineContacts(newContact.LinkedID)

			 if saveResult.RowsAffected > 0 {
				context.JSON(http.StatusOK,gin.H{"contact":consolidatedContact})
			 } else {
				context.JSON(http.StatusBadRequest,gin.H{"status":"RECORD NOT SAVED" ,"error":saveResult.Error})
			 }
		} else {
			 // no match in primary contacts  , proceed to secondary
			 //need to create the new contact since  there is no contact with this phone/email

			 newContact := models.Contact{
				PhoneNumber: userDetails.PhoneNumber,
				Email:userDetails.Email,
				LinkPrecedence : "primary",
		 	}

			saveResult := database.DB.Save(&newContact)

			consolidatedContact:=combineContacts(newContact.ID)

			if saveResult.RowsAffected > 0 {
				context.JSON(http.StatusOK,gin.H{"contact":consolidatedContact})
			} else {
				context.JSON(http.StatusBadRequest,gin.H{"status":"RECORD NOT SAVED" ,"error":saveResult.Error})
			}
			// none of phone/email matched to primary contact , create a new primary contact 
		}
}
}