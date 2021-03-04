package models

import (
	"go-dwh-api/app"
	u "go-dwh-api/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Statement of billing for the owner
type Statement struct {
	gorm.Model  `json:"-"`
	ContactID   uint32      `json:"-"`
	DueDate     time.Time   `json:"dueDate" gorm:"index:,sort:desc`
	Balance     float64     `json:"balance"`
	Assessments *Assessment `json:"assessment,omitempty" gorm:"foreignKey:StatementID"`
	PastDue     float64     `json:"pastDue"`
	Monthly     Monthly     `json:"monthlyStatement" gorm:"foreignKey:StatementID"`
}

// Assessment is different then a one time Charge as it is possible to be paid in installments.
//
type Assessment struct {
	gorm.Model  `json:"-"`
	StatementID uint32  `json:"-"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Balance     float64 `json:"balance"`
}

// Monthly is an archive of the last year of statements by month.
type Monthly struct {
	gorm.Model  `json:"-"`
	StatementID uint `json:"-"`
	Month       Month
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

// CreateStatement validates and create contact and sends the json
func (statement *Statement) CreateStatement() map[string]interface{} {

	app.GetDB().Create(statement)
	resp := u.Message(true, "Statement successfully created.")
	resp["statement"] = statement
	return resp
}

// GetCurrentStatement of the account number.
func GetCurrentStatement(contactID uint32) *Statement {

	statement := &Statement{}
	err := app.GetDB().Where("contact_id = ?", contactID).Preload(clause.Associations).First(&statement).Error
	if err != nil {
		u.Log.Error(err)
		return nil
	}
	return statement
}

// GetStatementHistory of all contacts
func GetStatementHistory(contactID uint32) []*Statement {
	statements := make([]*Statement, 0)

	err := app.GetDB().Where("contact_id = ?", contactID).Preload(clause.Associations).Find(&statements).Error
	if err != nil {
		u.Log.Error(err)
		return nil
	}

	return statements
}
