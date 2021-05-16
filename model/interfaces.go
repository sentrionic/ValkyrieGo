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

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(id string) (*User, error)
	Create(u *User) error
	FindByEmail(email string) (*User, error)
	Update(u *User) error
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
}
