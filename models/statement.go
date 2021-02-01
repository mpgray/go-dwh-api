package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Statement of billing for the owner
type Statement struct {
	gorm.Model
	DueDate time.Time `json:"dueDate"`
	Owed    float64   `json:"owed"`
	PastDue float64   `json:"pastDue"`
}

// GetCurrentStatement of the account number.
func GetCurrentStatement(id uint) *Contact {

	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		return nil
	}
	return contact
}

// GetStatementHistory of the account number for the past year.
func GetStatementHistory(user uint) []*Contact {

	contacts := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return contacts
}
