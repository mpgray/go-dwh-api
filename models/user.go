package models

import (
	"go-dwh-api/app"
	u "go-dwh-api/utils"

	e "github.com/dchest/validator"
	p "github.com/go-passwd/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//User is the basic login information of the user.
type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Validate incoming user details...
func (user *User) validate() (map[string]interface{}, bool) {

	if !e.IsValidEmail(user.Email) {
		u.Log.Warnf("User tried to create an account with an invalid email: %s", user.Email)
		return u.Message(false, "Email address is Invaild"), false
	}

	passwordValidator := p.Validator{p.MinLength(6, nil), p.MaxLength(40, nil), p.CommonPassword(nil)}
	passwordMessage := passwordValidator.Validate(user.Password)
	if passwordMessage != nil {
		u.Log.Warn(passwordMessage.Error())
		return u.Message(false, passwordMessage.Error()), false
	}

	//Email must be unique
	temp := &User{}

	//check for errors and duplicate emails
	err := app.GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		u.Log.Error("DB Connection Failed: While checking for duplicate emails")
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		u.Log.Warnf("User tried to register with an email that already exists: %s", user.Email)
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

// CreateUser the user's account
func (user *User) CreateUser() map[string]interface{} {

	if resp, ok := user.validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	app.GetDB().Create(user)

	if user.ID <= 0 {
		u.Log.Error("DB Connection Failed: While creating account")
		return u.Message(false, "Failed to create account, connection error.")
	}

	user.Password = "" //delete password

	u.Log.Infof("New Account-- Email: %s", user.Email)
	response := u.Message(true, "Account has been created")
	response["user"] = user
	return response
}

// GetUser returns nil when user not found in the database and the information it does
func GetUser(u uint64) *User {

	acc := &User{}
	app.GetDB().Table("user").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
