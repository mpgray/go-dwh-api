package models

import (
	"go-dwh-api/app"

	"github.com/gin-gonic/gin"
)

// Todo is a test struct
type Todo struct {
	UserID uint32 `json:"user_id"`
	Title  string `json:"title"`
}

// AccessDetails associates the UserId with the AccessUuid
type AccessDetails struct {
	AccessUUID string
	UserID     uint32
}

// FetchAuthenticatedID gets the authenticated id
func FetchAuthenticatedID(c *gin.Context, j interface{}) (uint32, error) {
	if err := c.ShouldBindJSON(&j); err != nil {
		app.UnprocessableEntityError(c, "Json Wrong during authentication")
		return 0, err
	}

	metadata, err := ExtractTokenMetadata(c.Request)
	if err != nil {

		return 0, err
	}
	return metadata.UserID, err
}
