package model

// Friend represents the api response of a user's friend.
type Friend struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	IsOnline bool   `json:"isOnline"`
}
