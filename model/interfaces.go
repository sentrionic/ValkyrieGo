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
	CreateDefaultChannel(c *Channel) error
	GenerateInviteLink(ctx context.Context, guildId string, isPermanent bool) (string, error)
	UpdateGuild(g *Guild) error
	GetGuildIdFromInvite(ctx context.Context, token string) (string, error)
	GetDefaultChannel(guildId string) (*Channel, error)
	InvalidateInvites(ctx context.Context, guild *Guild)
	RemoveMember(userId string, guildId string) error
	DeleteGuild(guildId string) error
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
}

type ChannelRepository interface {
	Create(c *Channel) error
	GetGuildDefault(guildId string) (*Channel, error)
}

type ImageRepository interface {
	UploadAvatar(header *multipart.FileHeader, directory string) (string, error)
	UploadImage(header *multipart.FileHeader, directory string) (string, error)
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
