package models

import (
	"fmt"
	u "go-hoa-api/utils"

	"github.com/jinzhu/gorm"
)

// Contact gets the Name, phone and UserID of the contact
type Contact struct {
	gorm.Model
	UserID     uint      `json:"user_id"` //The user that this contact belongs to
	Name       FullName  `json:"name"`
	Address    Address   `json:"address"`
	Properties []Address `json:"propertyAddresses"`
	Phone      string    `json:"phone"`
}

//FullName contains owners first middle and last name
type FullName struct {
	First  string `json:"first"`
	Middle string `json:"middle"`
	Last   string `json:"last"`
}

//Address contains all of the address information
type Address struct {
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	Unit  string `json:"unit"`
	City  string `json:"city"`
	State string `json:"state"`
	Zip   string `json:"zip"`
	Zip4  string `json:"zip4"`
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (contact *Contact) Validate() (map[string]interface{}, bool) {

	if contact.Name.First == "" {
		return u.Message(false, "Contact name should be on the payload"), false
	}

	if contact.Phone == "" {
		return u.Message(false, "Phone number should be on the payload"), false
	}

	if contact.UserID <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	//All the required parameters are present
	return u.Message(true, "success"), true
}

// Create makes the contact
func (contact *Contact) Create() map[string]interface{} {

	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	GetDB().Create(contact)

	resp := u.Message(true, "success")
	resp["contact"] = contact
	return resp
}

// GetContact by id
func GetContact(id uint) *Contact {

	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		return nil
	}
	return contact
}

// GetContacts gets all the contacts associated with the user.
func GetContacts(user uint) []*Contact {

	contacts := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return contacts
}
