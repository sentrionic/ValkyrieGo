package ws

// Subscribed Messages
const (
	JoinUserAction        = "joinUser"
	JoinGuildAction       = "joinGuild"
	JoinChannelAction     = "joinChannel"
	LeaveGuildAction      = "leaveGuild"
	LeaveRoomAction       = "leaveRoom"
	StartTypingAction     = "startTyping"
	StopTypingAction      = "stopTyping"
	ToggleOnlineAction    = "toggleOnline"
	ToggleOfflineAction   = "toggleOffline"
	GetRequestCountAction = "getRequestCount"
)

// Emitted Messages
const (
	NewMessageAction        = "new_message"
	EditMessageAction       = "edit_message"
	DeleteMessageAction     = "delete_message"
	AddChannelAction        = "add_channel"
	EditChannelAction       = "edit_channel"
	DeleteChannelAction     = "delete_channel"
	RemoveFromGuildAction   = "remove_from_guild"
	AddMemberAction         = "add_member"
	RemoveMemberAction      = "remove_member"
	NewDMNotificationAction = "new_dm_notification"
	NewNotificationAction   = "new_notification"
	ToggleOnlineEmission    = "toggle_online"
	ToggleOnlineOffline     = "toggle_offline"
	AddToTypingAction       = "addToTyping"
	RemoveFromTypingAction  = "removeFromTyping"
	SendRequestAction       = "send_request"
	AddRequestAction        = "add_request"
	AddFriendAction         = "add_friend"
	RemoveFriendAction      = "remove_friend"
)
