package repository

import (
	"database/sql"
	"fmt"
	"github.com/sentrionic/valkyrie/model"
	"gorm.io/gorm"
	"time"
)

// messageRepository is data/repository implementation
// of service layer MessageRepository
type messageRepository struct {
	DB *gorm.DB
}

// NewMessageRepository is a factory for initializing Message Repositories
func NewMessageRepository(db *gorm.DB) model.MessageRepository {
	return &messageRepository{
		DB: db,
	}
}

type MessageQuery struct {
	Id            string
	Text          *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FileType      *string
	Url           *string
	Filename      *string
	AttachmentId  *string
	UserId        string
	UserCreatedAt time.Time
	UserUpdatedAt time.Time
	Username      string
	Image         string
	IsOnline      bool
	Nickname      *string
	Color         *string
	IsFriend      bool
}

func (r *messageRepository) GetMessages(userId string, channel *model.Channel, cursor string) (*[]model.MessageResponse, error) {
	var result []MessageQuery

	memberSelect := ""
	memberJoin := ""
	memberWhere := ""
	if !channel.IsDM {
		memberSelect = "member.nickname, member.color,"
		memberJoin = "LEFT JOIN members member on messages.user_id = member.user_id"
		memberWhere = fmt.Sprintf("AND member.guild_id = %s::text", *channel.GuildID)
	}

	crs := ""
	if cursor != "" {
		date := cursor[:len(cursor)-6]
		crs = fmt.Sprintf("AND messages.created_at < '%s'", date)
	}

	err := r.DB.
		Raw(fmt.Sprintf(`
		SELECT messages.id,
			messages.text,
			messages.created_at,
			messages.updated_at,
			a.file_type,
			a.url,
			a.filename,
			a.id                as "attachment_id",
			users.id         as "user_id",
			users.created_at as "user_created_at",
			users.updated_at as "user_updated_at",
			users.username,
			users.image,
			users.is_online,
			%s 
			EXISTS(
			  SELECT 1
			  FROM users
			   LEFT JOIN friends f ON users.id = f.user_id
			  WHERE f.friend_id = messages.user_id
				AND f.user_id = @userId) as "isFriend"
		FROM messages
		LEFT JOIN "users"
		ON users.id = messages.user_id
		LEFT JOIN attachments a
		ON a.message_id = messages.id
		%s
		WHERE messages.channel_id = @channelId
		%s 
		%s 
		ORDER BY messages.created_at DESC
		LIMIT 35
`, memberSelect, memberJoin, memberWhere, crs),
			sql.Named("userId", userId),
			sql.Named("channelId", channel.ID)).
		Scan(&result).Error

	var messages []model.MessageResponse

	for _, m := range result {

		var attachment *model.Attachment = nil
		if m.AttachmentId != nil {
			attachment = &model.Attachment{
				Url:      *m.Url,
				FileType: *m.FileType,
				Filename: *m.Filename,
			}
		}

		message := model.MessageResponse{
			Id:         m.Id,
			Text:       m.Text,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
			Attachment: attachment,
			User: model.MemberResponse{
				Id:        m.UserId,
				Username:  m.Username,
				Image:     m.Image,
				IsOnline:  m.IsOnline,
				CreatedAt: m.UserCreatedAt,
				UpdatedAt: m.UserUpdatedAt,
				Nickname:  m.Nickname,
				Color:     m.Color,
				IsFriend:  m.IsFriend,
			},
		}
		messages = append(messages, message)
	}

	return &messages, err
}

func (r *messageRepository) CreateMessage(message *model.Message) error {
	return r.DB.Create(&message).Error
}

func (r *messageRepository) UpdateMessage(message *model.Message) error {
	return r.DB.Save(&message).Error
}

func (r *messageRepository) DeleteMessage(message *model.Message) error {
	return r.DB.Delete(message).Error
}

func (r *messageRepository) GetById(messageId string) (*model.Message, error) {
	var message model.Message
	err := r.DB.Where("id = ?", messageId).First(&message).Error
	return &message, err
}
