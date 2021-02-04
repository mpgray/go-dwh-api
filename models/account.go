package models

import (
	u "go-hoa-api/utils"
	"os"

	e "github.com/dchest/validator" //email validation
	"github.com/dgrijalva/jwt-go"
	p "github.com/go-passwd/validator" // password validation
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Token is a struct which gets the jwt for our claim.
type Token struct {
	UserID uint
	jwt.StandardClaims
}

//Account is used to login the user with their email and password using tokens from jwt
type Account struct {
	gorm.Model
	Email    string   `json:"email"`
	Password string   `json:"password"`
	UserType UserType `json:"userType"`
	Token    string   `json:"token" sql:"-"`
}

// UserType represents what the use could be. Home Owner, manager, SuperUser
type UserType uint8

const (
	// OWNER is the Home Owner
	OWNER UserType = iota
	// MANAGER is the HOA
	MANAGER
	// SUPERUSER is administrator of this application
	SUPERUSER
)

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !e.IsValidEmail(account.Email) {
		return u.Message(false, "Email address is required"), false
	}

	passwordValidator := p.Validator{p.MinLength(6, nil), p.MaxLength(40, nil), p.CommonPassword(nil)}
	passwordMessage := passwordValidator.Validate(account.Password)
	if passwordMessage != nil {
		return u.Message(false, passwordMessage.Error()), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

// Create the user's account
func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

// Login to the account using bcrypt and JWT
func Login(email, password string) map[string]interface{} {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

// GetUser returns nil when user not found in the database and the information it does
func GetUser(u uint) *Account {

	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}