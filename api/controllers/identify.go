package controllers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	database "github.com/dg222599/BiteSpeed-Backend-Assignment/database"
	models "github.com/dg222599/BiteSpeed-Backend-Assignment/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	 PhoneNumber interface{}
	 Email interface{}
}

type combinedContact struct {
	PrimaryContactId uint `json:"primaryContactId"`
	Emails []*string	`json:"emails"`
	PhoneNumbers []*string `json:"phoneNumbers"`
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
	
	completeContactData := combinedContact{}

	//assign ID
	uniqueEmails := make(map[string]bool)
	uniquePhoneNumbers := make(map[string]bool)

	completeContactData.PrimaryContactId = primaryID

	var primaryContact models.Contact
	primaryContactResult := database.DB.Where("id=?",primaryID).Find(&primaryContact)
	if primaryContactResult.RowsAffected <=0 {
		 return completeContactData
	}

	//assign first email
	
	if primaryContact.Email != nil {
		completeContactData.Emails = append(completeContactData.Emails,primaryContact.Email)
		uniqueEmails[*primaryContact.Email] = true
	}

	if primaryContact.PhoneNumber != nil {
		completeContactData.PhoneNumbers = append(completeContactData.PhoneNumbers, primaryContact.PhoneNumber)
		uniquePhoneNumbers[*primaryContact.PhoneNumber] = true
	}
	
	

	// find all the contacts which are related to this primary contact
	var relatedContacts []models.Contact
	result:= database.DB.Where("linked_id=?",primaryID).Order("created_at").Find(&relatedContacts)

	if result.RowsAffected <=0 {
		 return completeContactData
	}


	for index:=0;index<len(relatedContacts);index++ {
		 
		// need to append email,phonenumber,id in seconddarycontactIDs

		var newContact bool
		newContact  = false
		if relatedContacts[index].Email != nil && !uniqueEmails[*relatedContacts[index].Email] {
			completeContactData.Emails = append(completeContactData.Emails,relatedContacts[index].Email)
			uniqueEmails[*relatedContacts[index].Email] = true
			newContact = true
		}

		if relatedContacts[index].PhoneNumber != nil && !uniquePhoneNumbers[*relatedContacts[index].PhoneNumber] {
			completeContactData.PhoneNumbers = append(completeContactData.PhoneNumbers,relatedContacts[index].PhoneNumber)
			uniquePhoneNumbers[*relatedContacts[index].PhoneNumber] = true
			newContact = true
		}
		
		// completeContactData.PhoneNumbers = append(completeContactData.PhoneNumbers,relatedContacts[index].PhoneNumber)
		if newContact {
			completeContactData.SecondaryContactIds = append(completeContactData.SecondaryContactIds, relatedContacts[index].ID)
		}
		
	}

	return completeContactData
}
func ValidateRequest(userDetails *User) bool {
	// at least one non null value and should be string

	
	
	if  userDetails.Email == nil && userDetails.PhoneNumber == nil {
		 return false
	} else if (userDetails.Email == nil || userDetails.PhoneNumber == nil) {
		
		if((userDetails.Email!=nil) && (reflect.TypeOf(userDetails.Email).Kind()==reflect.String) && userDetails.Email!=""){
			 return true
		} else if ((userDetails.PhoneNumber!=nil) && (reflect.TypeOf(userDetails.PhoneNumber).Kind()==reflect.String) && userDetails.PhoneNumber!=""){
			 return true
		} else {
			 return false
		}
 		
	} else {
		
		  if(reflect.TypeOf(userDetails.Email).Kind()==reflect.String  && reflect.TypeOf(userDetails.PhoneNumber).Kind()==reflect.String && userDetails.PhoneNumber!="" && userDetails.Email!=""){
			 return true
		  }  else {
			 return false
		  }
	}
}
func LinkIdentity(context *gin.Context){
   
	 primaryPrecedence := "primary"
	 secondaryPrecedence := "secondary"
	 var userDetails User

	 
	 if err:= context.BindJSON(&userDetails) ; err!=nil{
		context.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return 
	 }

	 validationResult:=ValidateRequest(&userDetails)

     
	  
	 if !validationResult {
		 //both the fields are empty
		 exampleRequest := User{
			Email:"abc@gmail.com",
			PhoneNumber:"12345678",
		 }
		 context.JSON(http.StatusBadRequest,gin.H{"message":"Bad-Request - all the non null fields should be a non-empty string","example":exampleRequest})
		 return
	 }

	 

	 var alreadyPresentRecord models.Contact


	 var contactAlreadyPresent *gorm.DB

	 if userDetails.Email != nil && userDetails.PhoneNumber != nil {
		contactAlreadyPresent = database.DB.Where("phone_number=? and email=?",userDetails.PhoneNumber,userDetails.Email).First(&alreadyPresentRecord)
	 } else if userDetails.Email != nil {
		contactAlreadyPresent = database.DB.Where("email=?",userDetails.Email).First(&alreadyPresentRecord)
	 } else if userDetails.PhoneNumber != nil {
	    contactAlreadyPresent = database.DB.Where("phone_number=?",userDetails.PhoneNumber).First(&alreadyPresentRecord)	 
	 }


	 if contactAlreadyPresent.RowsAffected > 0 {
		
		var connectionID uint
		if alreadyPresentRecord.LinkedID !=nil {
			 connectionID = *alreadyPresentRecord.LinkedID
		} else {
			 connectionID = alreadyPresentRecord.ID
		}
		consolidatedContact:=combineContacts(connectionID)
		context.JSON(http.StatusOK,gin.H{"contact":consolidatedContact})
		return
	 } else {
		
	    var primaryEmailContact,primaryPhoneContact models.Contact

		var phoneResultPrimary,emailResultPrimary * gorm.DB

		var secondaryRecordPresent *gorm.DB
		var bogus models.Contact

		if userDetails.PhoneNumber != nil {
			 secondaryRecordPresent = database.DB.Where("phone_number=? and link_precedence=?",userDetails.PhoneNumber,"secondary").Order("created_at").First(&bogus)

		}

		if secondaryRecordPresent.RowsAffected > 0 {
			 
			 consolidatedContact := combineContacts(*bogus.LinkedID)
			 context.JSON(http.StatusOK,gin.H{"contact":consolidatedContact})
			 return
		}

		if userDetails.Email != nil {
			secondaryRecordPresent = database.DB.Where("email=? and link_precedence=?",userDetails.Email,"secondary").Order("created_at").First(&bogus)

	    }

		if secondaryRecordPresent.RowsAffected > 0 {
			
			consolidatedContact := combineContacts(*bogus.LinkedID)

			context.JSON(http.StatusOK,gin.H{"contact":consolidatedContact})
			return
		}


		
		// denotes the record where a primary contact has same email/phone as current contact
		if userDetails.PhoneNumber != nil {
			phoneResultPrimary = database.DB.Where("phone_number=? and link_precedence=?",userDetails.PhoneNumber,"primary").Order("created_at").First(&primaryPhoneContact)
	 	}

		if userDetails.Email != nil {
			emailResultPrimary = database.DB.Where("email=? and link_precedence=?",userDetails.Email,"primary").Order("created_at").First(&primaryEmailContact)
		}
		
		
		if (emailResultPrimary!= nil && emailResultPrimary.RowsAffected>0 ) && (phoneResultPrimary != nil && phoneResultPrimary.RowsAffected > 0) {
			  //found primary match in both(can be same contact can be not)
			 // email coming from one contact and phone coming from another contact
			 //create link either way
			 var newestContact models.Contact
			 newestContactResult := database.DB.Where("id in ?",[]uint{primaryEmailContact.ID,primaryPhoneContact.ID}).Order("created_at desc").First(&newestContact)
			 
			 if newestContactResult.RowsAffected <=0 {
				 //
			 }
			 var updatedResult *gorm.DB
			 fmt.Println("------------------------->",newestContact.ID)
			 if newestContact.ID == primaryEmailContact.ID {
				 
				 updatedResult = database.DB.Model(&models.Contact{}).Where("id=?",newestContact.ID).Update("linked_id",primaryPhoneContact.ID)

			 } else {
				 
				 updatedResult = database.DB.Model(&models.Contact{}).Where("id=?",newestContact.ID).Update("linked_id",primaryEmailContact.ID)
				 
			 }
			 
			 updatedResult = database.DB.Model(&models.Contact{}).Where("id=?",newestContact.ID).Update("link_precedence","secondary")

			 

			 if updatedResult.RowsAffected > 0 {
				
				var consolidatedContact combinedContact
				if newestContact.ID == primaryEmailContact.ID {
					consolidatedContact = combineContacts(primaryPhoneContact.ID)
				} else {
					consolidatedContact = combineContacts(primaryEmailContact.ID) 
				}
				
				context.JSON(http.StatusOK,gin.H{"contact":consolidatedContact})
				return
			 }
		} else if (emailResultPrimary!= nil && emailResultPrimary.RowsAffected>0 ) || (phoneResultPrimary != nil && phoneResultPrimary.RowsAffected > 0){

			 //found primary match in only one
			 newContact := models.Contact{
				LinkPrecedence : &secondaryPrecedence,
				Email: nil,
				PhoneNumber: nil,
		 	 }

			if userDetails.Email != nil {
				if email, ok := userDetails.Email.(string); ok {
					// Type assertion successful, assign the string value to userEmail
					newContact.Email = &email
				} 
			} 
			
			// Check if userDetails.PhoneNumber is not nil before type assertion
			if userDetails.PhoneNumber != nil {
				if phoneNumber, ok := userDetails.PhoneNumber.(string); ok {
					// Type assertion successful, assign the string value to userPhoneNumber
					newContact.PhoneNumber = &phoneNumber
				}
			}	

			 
			 if emailResultPrimary.RowsAffected > 0 {
				 newContact.LinkedID = &primaryEmailContact.ID
			 } else {
				 newContact.LinkedID = &primaryPhoneContact.ID
			 }

			 saveResult := database.DB.Save(&newContact)
			 if saveResult.RowsAffected <=0 {
					return
			 }

			 
			 consolidatedContact:=combineContacts(*newContact.LinkedID)

			 if saveResult.RowsAffected > 0 {
				context.JSON(http.StatusOK,gin.H{"contact":consolidatedContact})
			 } else {
				context.JSON(http.StatusBadRequest,gin.H{"status":"RECORD NOT SAVED" ,"error":saveResult.Error})
			 }
		} else {
			 // no match in primary contacts  , proceed to secondary
			 //need to create the new contact since  there is no contact with this phone/email
			 fmt.Println("reaced till print 0")
			 newContact := models.Contact{
				LinkPrecedence : &primaryPrecedence,
				Email: nil,
				PhoneNumber: nil,
		 	 }
			 fmt.Println("reaced till print 1")
			if userDetails.Email != nil {
				if email, ok := userDetails.Email.(string); ok {
					// Type assertion successful, assign the string value to userEmail
					newContact.Email = &email
				} 
			} 
			
			// Check if userDetails.PhoneNumber is not nil before type assertion
			if userDetails.PhoneNumber != nil {
				if phoneNumber, ok := userDetails.PhoneNumber.(string); ok {
					// Type assertion successful, assign the string value to userPhoneNumber
					newContact.PhoneNumber = &phoneNumber
				}
			}
			 
			fmt.Println("reached till print -2")
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