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
func GetCurrentStatement(id uint) *Statement {

	statement := &Statement{}
	err := GetDB().Table("statements").Where("id = ?", id).First(statement).Error
	if err != nil {
		return nil
	}
	return statement
}

// GetStatementHistory of the account number for the past year.
func GetStatementHistory(user uint) []*Statement {

	statements := make([]*Statement, 0)
	err := GetDB().Table("statements").Where("user_id = ?", user).Find(&statements).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return statements
}
