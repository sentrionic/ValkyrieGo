package model

type FriendRequest struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	Type     int    `json:"type"`
}
