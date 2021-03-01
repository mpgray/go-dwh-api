package models

/* Contact Model
This model contains all the profile information, including billing,
of a contact. A contact is defined as either an association profile,
an owner, or a tenent (non owner occupant)
*/
import (
	"go-dwh-api/app"
	u "go-dwh-api/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Contact gets the Name, phone and UserID of the contact
type Contact struct {
	gorm.Model `json:"-"`
	ID         uint32       //This user's ID
	UserID     uint32       `json:"-"` //The user that this contact belongs to
	Name       *FullName    `gorm:"foreignkey:ContactID" json:"name,omitempty"`
	Address    *Address     `gorm:"foreignkey:ContactID" json:"address,omitempty"` // Mailing address
	Statement  *[]Statement `gorm:"foreignkey:ContactID" json:"statement,omitempty"`
	//Properties []Address `json:"propertyAddresses" gorm:"foreignKey:ID"` // The properties owned by the owner
	Phone *Phone `gorm:"foreignkey:ContactID" json:"phone,omitempty"`
}

const (
	// NAME contact sub category
	NAME string = "Name"
	// ADDRESS contact sub category
	ADDRESS string = "Address"
	// PHONE contact sub category
	PHONE string = "Phone"
)

//FullName contains owners first middle and last name
type FullName struct {
	gorm.Model `json:"-"`
	ContactID  uint32 `json:"-"`
	First      string `json:"first"`
	Middle     string `json:"middle,omitempty"`
	Last       string `json:"last"`
}

//Address contains all of the address information
type Address struct {
	gorm.Model `json:"-"`
	ContactID  uint32 `json:"-"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	Unit       string `json:"unit,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	Zip        string `json:"zip"`
	Zip4       string `json:"zip4,omitempty"`
}

// Phone Contains different phone numbers of the home owner
type Phone struct {
	gorm.Model `json:"-"`
	ContactID  uint32       `json:"-"`
	Cell       string       `json:"cellPhone,omitempty"`
	Home       string       `json:"homePhone,omitempty"`
	Work       string       `json:"workPhone,omitempty"`
	Other      string       `json:"otherPhone,omitempty"`
	Primary    PrimaryPhone `json:"primaryPhone"`
}

// ContactID is what you get back when a user wants a certain user
// also what the user must send to get a certain user { ID: 0 }
type ContactID struct {
	ID uint32 `json:"ID"`
}

// PrimaryPhone is an enum to determain primary phone contact
type PrimaryPhone uint8

const (
	// CELL is Primary phone
	CELL PrimaryPhone = iota
	// HOME is Primary phone
	HOME
	// WORK is Primary phone
	WORK
	// OTHER number is primary phone
	OTHER
)

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (contact *Contact) validate() (map[string]interface{}, bool) {
	message := "Missing required fields:\n"
	isValid := true

	if contact.Name.First == "" {
		message += "First Name\n"
		isValid = false
	}

	if contact.Name.Last == "" {
		message += "Last Name\n"
		isValid = false
	}

	if contact.Phone.Cell == "" {
		message += "Cell Phone\n"
		isValid = false
	}

	if contact.Address.Line1 == "" || contact.Address.City == "" || contact.Address.State == "" || contact.Address.Zip == "" {
		message += "Address, City, State or Zip\n"
		isValid = false
	}

	if contact.UserID <= 0 {
		message = "Unknow User" // if user isnt' know, no other data should display
		isValid = false
	}

	if isValid {
		message = "success"
	}

	//Shows what is missing and if successful, displays success message
	return u.Message(isValid, message), isValid
}

// CreateContact validates and create contact and sends the json
func (contact *Contact) CreateContact() map[string]interface{} {

	if resp, ok := contact.validate(); !ok {
		return resp
	}
	app.GetDB().Create(contact)
	resp := u.Message(true, "Contact successfully created.")
	resp["contact"] = contact
	return resp
}

// GetContact by id
func GetContact(contactID uint32, user uint32) *Contact {
	contact := &Contact{}

	err := app.GetDB().Where("user_id = ? and ID = ?", user, contactID).Preload(clause.Associations).First(&contact).Error
	if err != nil {
		u.Log.Error(err)
		return nil
	}
	return contact
}

// GetContacts gets all the contacts associated with the user.
func GetContacts(user uint32) []*Contact {

	contacts := make([]*Contact, 0)

	err := app.GetDB().Where("user_id = ?", user).Preload(clause.Associations).Find(&contacts).Error
	if err != nil {
		u.Log.Error(err)
	}

	return contacts
}

//SearchContacts gets the name and addresses of the contacts to search
func SearchContacts(user uint32) []*Contact {
	contacts := make([]*Contact, 0)

	err := app.GetDB().Where("user_id = ?", user).Joins(NAME).Joins(ADDRESS).Find(&contacts).Error
	if err != nil {
		u.Log.Error(err)
	}

	return contacts
}

//GetAllContactInfo takes the category and user id and returns that users contacts' category, for example phone numbers
func GetAllContactInfo(user uint32, category string) []*Contact {
	contacts := make([]*Contact, 0)

	err := app.GetDB().Where("user_id = ?", user).Joins(category).Find(&contacts).Error
	if err != nil {
		u.Log.Error(err)
	}

	return contacts
}

// GetContactInfo gets the sub category of a contact. For example, phone number.
func GetContactInfo(user uint32, contactID uint32, category string) *Contact {
	contact := &Contact{}
	err := app.GetDB().Where("user_id = ? and contact_id = ?", user, contactID).Joins(category).First(&contact).Error
	if err != nil {
		u.Log.Error(err)
		return nil
	}
	return contact
}
