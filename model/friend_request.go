package model

//FriendRequest contains all info to display request.
//Type stands for the type of the request.
//1: Incoming,
//0: Outgoing
type FriendRequest struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	// 1: Incoming, 0: Outgoing
	Type int `json:"type" enums:"0,1"`
} //@name FriendRequest
