package model

type DirectMessage struct {
	Id   string `json:"id"`
	User DMUser `json:"user"`
}

type DMUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	IsOnline bool   `json:"isOnline"`
	IsFriend bool   `json:"IsFriend"`
}
