package models

import (
	"go-dwh-api/app"
	u "go-dwh-api/utils"

	e "github.com/dchest/validator"
	p "github.com/go-passwd/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//Account is used to login the user with their email and password using tokens from jwt
type Account struct {
	gorm.Model
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	UserType  UserType `json:"userType"`
	ManagerID *uint
	Owners    []Account `gorm:"foreignkey:ManagerID"`
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

/*

func createToken(account *Account) string {
	//Create new JWT token for the newly registered account

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = &Token{UserID: account.ID}
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(os.Getenv("token_password")))

	return token
}

// Login to the account using bcrypt and JWT
func Login(email, password string) map[string]interface{} {

	account := &Account{}
	err := app.GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warn("Attempt to login resulted in email not found.")
			return u.Message(false, "Email address not found")
		}
		u.Log.Error("DB Connection Failed: During login attempt")
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		u.Log.Warn("Attempt to login with invalid credentials: %s ")
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
} */

// GetUser returns nil when user not found in the database and the information it does
func GetUser(u uint) *Account {

	acc := &Account{}
	app.GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
