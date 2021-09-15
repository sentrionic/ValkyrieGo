package model

import (
	"context"
	"mime/multipart"
)

// GuildService defines methods related to guild operations the handler layer expects
// any service it interacts with to implement
type GuildService interface {
	GetUser(uid string) (*User, error)
	GetGuild(id string) (*Guild, error)
	GetUserGuilds(uid string) (*[]GuildResponse, error)
	GetGuildMembers(userId string, guildId string) (*[]MemberResponse, error)
	CreateGuild(guild *Guild) error
	GenerateInviteLink(ctx context.Context, guildId string, isPermanent bool) (string, error)
	UpdateGuild(guild *Guild) error
	GetGuildIdFromInvite(ctx context.Context, token string) (string, error)
	GetDefaultChannel(guildId string) (*Channel, error)
	InvalidateInvites(ctx context.Context, guild *Guild)
	RemoveMember(userId string, guildId string) error
	UnbanMember(userId string, guildId string) error
	DeleteGuild(guildId string) error
	GetBanList(guildId string) (*[]BanResponse, error)
	GetMemberSettings(userId string, guildId string) (*MemberSettings, error)
	UpdateMemberSettings(settings *MemberSettings, userId string, guildId string) error
	FindUsersByIds(ids []string, guildId string) (*[]User, error)
	UpdateMemberLastSeen(userId, guildId string) error
}

// MessageService defines methods related to message operations the handler layer expects
// any service it interacts with to implement
type MessageService interface {
	GetMessages(userId string, channel *Channel, cursor string) (*[]MessageResponse, error)
	CreateMessage(message *Message) error
	UpdateMessage(message *Message) error
	DeleteMessage(message *Message) error
	UploadFile(header *multipart.FileHeader, channelId string) (*Attachment, error)
	Get(messageId string) (*Message, error)
}

// GuildRepository defines methods related to guild db operations the service layer expects
// any repository it interacts with to implement
type GuildRepository interface {
	FindUserByID(uid string) (*User, error)
	FindByID(id string) (*Guild, error)
	List(uid string) (*[]GuildResponse, error)
	GuildMembers(userId string, guildId string) (*[]MemberResponse, error)
	Create(guild *Guild) error
	Save(guild *Guild) error
	RemoveMember(userId string, guildId string) error
	Delete(guildId string) error
	UnbanMember(userId string, guildId string) error
	GetBanList(guildId string) (*[]BanResponse, error)
	GetMemberSettings(userId string, guildId string) (*MemberSettings, error)
	UpdateMemberSettings(settings *MemberSettings, userId string, guildId string) error
	FindUsersByIds(ids []string, guildId string) (*[]User, error)
	GetMember(userId, guildId string) (*User, error)
	UpdateMemberLastSeen(userId, guildId string) error
	GetMemberIds(guildId string) (*[]string, error)
}

// FileRepository defines methods related to file upload the service layer expects
// any repository it interacts with to implement
type FileRepository interface {
	UploadAvatar(header *multipart.FileHeader, directory string) (string, error)
	UploadFile(header *multipart.FileHeader, directory, filename, mimetype string) (string, error)
	DeleteImage(key string) error
}

// MailRepository defines methods related to mail operations the service layer expects
// any repository it interacts with to implement
type MailRepository interface {
	SendResetMail(email string, html string) error
}

// RedisRepository defines methods related to the redis db the service layer expects
// any repository it interacts with to implement
type RedisRepository interface {
	SetResetToken(ctx context.Context, id string) (string, error)
	GetIdFromToken(ctx context.Context, token string) (string, error)
	SaveInvite(ctx context.Context, guildId string, id string, isPermanent bool) error
	GetInvite(ctx context.Context, token string) (string, error)
	InvalidateInvites(ctx context.Context, guild *Guild)
}

// MessageRepository defines methods related message db operations the service layer expects
// any repository it interacts with to implement
type MessageRepository interface {
	GetMessages(userId string, channel *Channel, cursor string) (*[]MessageResponse, error)
	CreateMessage(message *Message) error
	UpdateMessage(message *Message) error
	DeleteMessage(message *Message) error
	GetById(messageId string) (*Message, error)
}

// SocketService defines methods related emitting websocket events the service layer expects
// any repository it interacts with to implement
type SocketService interface {
	EmitNewMessage(room string, message *MessageResponse)
	EmitEditMessage(room string, message *MessageResponse)
	EmitDeleteMessage(room, messageId string)

	EmitNewChannel(room string, channel *ChannelResponse)
	EmitEditChannel(room string, channel *ChannelResponse)
	EmitDeleteChannel(channel *Channel)

	EmitEditGuild(guild *Guild)
	EmitDeleteGuild(guildId string, members []string)
	EmitRemoveFromGuild(memberId, guildId string)

	EmitAddMember(room string, member *User)
	EmitRemoveMember(room, memberId string)

	EmitNewDMNotification(channelId string, user *User)
	EmitNewNotification(guildId, channelId string)

	EmitSendRequest(room string)
	EmitAddFriendRequest(room string, request *FriendRequest)
	EmitAddFriend(user, member *User)
	EmitRemoveFriend(userId, memberId string)
}
