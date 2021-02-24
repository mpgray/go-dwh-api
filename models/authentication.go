package models

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

//TODO: Implement this so that a user_id isn't sent to token just an accessUUID that can get the userID
/*
// FetchAuth gets the userID with the AccessUUID
func FetchAuth(authD *AccessDetails) (uint32, error) {
	userid, err := app.GetRedis().Get(authD.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID64, _ := strconv.ParseUint(userid, 10, 64)
	userID := uint32(userID64)
	return userID, err
} */
