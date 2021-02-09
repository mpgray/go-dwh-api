package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Statement of billing for the owner
type Statement struct {
	gorm.Model
	ID          uint
	DueDate     time.Time  `json:"dueDate"`
	Balance     float64    `json:"balance"`
	Assessments Assessment `json:"assessment" gorm:"foreignKey:ID"`
	PastDue     float64    `json:"pastDue"`
	Monthly     Monthly    `json:"monthlyStatement"  gorm:"ID"`
}

// Assessment is different then a one time Charge as it is possible to be paid in installments.
//
type Assessment struct {
	gorm.Model
	ID      uint
	Name    string  `json:"name"`
	Amount  float64 `json:"amount"`
	Balance float64 `json:"balance"`
}

// Monthly is an archive of the last year of statements by month.
type Monthly struct {
	ID    uint
	Month Month
}

// Month is  an enum with 0 current month and actual months
type Month uint8

//Enum of months with current month being 0
const (
	CURRENT Month = iota
	JANUARY
	FEBRUARY
	MARCH
	APRIL
	MAY
	JUNE
	JULY
	AUGUST
	SEPTEMBER
	OCTOBER
	NOVEMBER
	DECEMBER
)

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
