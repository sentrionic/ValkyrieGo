package model

/*
	FriendRequest contains all info to display request.
	Type stands for the type of the request.
	1: Incoming,
	0: Outgoing
*/
type FriendRequest struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	Type     int    `json:"type"`
}
