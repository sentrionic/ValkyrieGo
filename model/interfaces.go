package model

import (
	"context"
	"mime/multipart"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type UserService interface {
	Get(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Register(u *User) error
	Login(u *User) error
	UpdateAccount(u *User) error
	CheckEmail(email string) bool
	ChangeAvatar(header *multipart.FileHeader, directory string) (string, error)
	DeleteImage(key string) error
	ChangePassword(password string, u *User) error
	ForgotPassword(ctx context.Context, u *User) error
	ResetPassword(ctx context.Context, password string, token string) (*User, error)
}

type FriendService interface {
	GetFriends(id string) (*[]Friend, error)
	GetRequests(id string) (*[]FriendRequest, error)
	GetMemberById(id string) (*User, error)
	DeleteRequest(memberId string, userId string) error
	RemoveFriend(memberId string, userId string) error
	SaveRequests(user *User) error
}

type GuildService interface {
	GetUser(uid string) (*User, error)
	GetGuild(id string) (*Guild, error)
	GetUserGuilds(uid string) (*[]GuildResponse, error)
	GetGuildMembers(userId string, guildId string) (*[]MemberResponse, error)
	CreateGuild(g *Guild) error
	GenerateInviteLink(ctx context.Context, guildId string, isPermanent bool) (string, error)
	UpdateGuild(g *Guild) error
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

type ChannelService interface {
	CreateChannel(channel *Channel) error
	GetChannels(userId string, guildId string) (*[]ChannelResponse, error)
	Get(channelId string) (*Channel, error)
	GetPrivateChannelMembers(channelId string) (*[]string, error)
	GetDirectMessages(userId string) (*[]DirectMessage, error)
	GetDirectMessageChannel(userId string, memberId string) (*string, error)
	GetDMByUserAndChannel(userId string, channelId string) (string, error)
	AddDMChannelMembers(memberIds []string, channelId string, userId string) error
	SetDirectMessageStatus(dmId string, userId string, isOpen bool) error
	DeleteChannel(channel *Channel) error
	UpdateChannel(channel *Channel) error
	CleanPCMembers(channelId string) error
	AddPrivateChannelMembers(memberIds []string, channelId string) error
	RemovePrivateChannelMembers(memberIds []string, channelId string) error
	IsChannelMember(channel *Channel, userId string) error
	OpenDMForAll(dmId string) error
}

type MessageService interface {
	GetMessages(userId string, channel *Channel, cursor string) (*[]MessageResponse, error)
	CreateMessage(message *Message) error
	UpdateMessage(message *Message) error
	DeleteMessage(message *Message) error
	UploadFile(header *multipart.FileHeader, channelId string) (*Attachment, error)
	Get(messageId string) (*Message, error)
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(id string) (*User, error)
	Create(u *User) error
	FindByEmail(email string) (*User, error)
	Update(u *User) error
}

type FriendRepository interface {
	FindByID(id string) (*User, error)
	FriendsList(id string) (*[]Friend, error)
	RequestList(id string) (*[]FriendRequest, error)
	DeleteRequest(memberId string, userId string) error
	RemoveFriend(memberId string, userId string) error
	Save(user *User) error
}

type GuildRepository interface {
	FindUserByID(uid string) (*User, error)
	FindByID(id string) (*Guild, error)
	List(uid string) (*[]GuildResponse, error)
	GuildMembers(userId string, guildId string) (*[]MemberResponse, error)
	Create(g *Guild) error
	Save(g *Guild) error
	RemoveMember(userId string, guildId string) error
	Delete(guildId string) error
	UnbanMember(userId string, guildId string) error
	GetBanList(guildId string) (*[]BanResponse, error)
	GetMemberSettings(userId string, guildId string) (*MemberSettings, error)
	UpdateMemberSettings(settings *MemberSettings, userId string, guildId string) error
	FindUsersByIds(ids []string, guildId string) (*[]User, error)
	GetMember(userId, guildId string) (*User, error)
	UpdateMemberLastSeen(userId, guildId string) error
}

type ChannelRepository interface {
	Create(c *Channel) error
	GetGuildDefault(guildId string) (*Channel, error)
	Get(userId string, guildId string) (*[]ChannelResponse, error)
	GetDirectMessages(userId string) (*[]DirectMessage, error)
	GetDirectMessageChannel(userId string, memberId string) (*string, error)
	GetById(channelId string) (*Channel, error)
	GetPrivateChannelMembers(channelId string) (*[]string, error)
	AddDMChannelMembers(members []DMMember) error
	SetDirectMessageStatus(dmId string, userId string, isOpen bool) error
	DeleteChannel(channel *Channel) error
	UpdateChannel(channel *Channel) error
	CleanPCMembers(channelId string) error
	AddPrivateChannelMembers(memberIds []string, channelId string) error
	RemovePrivateChannelMembers(memberIds []string, channelId string) error
	FindDMByUserAndChannelId(channelId, userId string) (string, error)
	OpenDMForAll(dmId string) error
}

type FileRepository interface {
	UploadAvatar(header *multipart.FileHeader, directory string) (string, error)
	UploadFile(header *multipart.FileHeader, directory, filename, mimetype string) (string, error)
	DeleteImage(key string) error
}

type MailRepository interface {
	SendMail(email string, html string) error
}

type RedisRepository interface {
	SetResetToken(ctx context.Context, id string) (string, error)
	GetIdFromToken(ctx context.Context, token string) (string, error)
	SaveInvite(ctx context.Context, guildId string, id string, isPermanent bool) error
	GetInvite(ctx context.Context, token string) (string, error)
	InvalidateInvites(ctx context.Context, guild *Guild)
}

type MessageRepository interface {
	GetMessages(userId string, channel *Channel, cursor string) (*[]MessageResponse, error)
	CreateMessage(message *Message) error
	UpdateMessage(message *Message) error
	DeleteMessage(message *Message) error
	GetById(messageId string) (*Message, error)
}
