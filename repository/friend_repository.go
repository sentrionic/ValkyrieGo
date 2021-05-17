package repository

import (
	"database/sql"
	"errors"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"gorm.io/gorm"
)

// friendRepository is data/repository implementation
// of service layer FriendRepository
type friendRepository struct {
	DB *gorm.DB
}

// NewFriendRepository is a factory for initializing User Repositories
func NewFriendRepository(db *gorm.DB) model.FriendRepository {
	return &friendRepository{
		DB: db,
	}
}

func (r *friendRepository) FriendsList(id string) (*[]model.Friend, error) {
	var u []model.Friend

	result := r.DB.
		Table("users").
		Joins(`JOIN friends ON friends.user_id = "users".id`).
		Where("friends.friend_id = ?", id).
		Find(&u)

	return &u, result.Error
}

func (r *friendRepository) RequestList(id string) (*[]model.FriendRequest, error) {
	var fr []model.FriendRequest

	result := r.DB.
		Raw(`
		  select u.id, u.username, u.image, 1 as "type" from users u
		  join friend_requests fr on u.id = fr."sender_id"
		  where fr."receiver_id" = @id
		  UNION
		  select u.id, u.username, u.image, 0 as "type" from users u
		  join friend_requests fr on u.id = fr."receiver_id"
		  where fr."sender_id" = @id
		  order by username;
		`, sql.Named("id", id)).Find(&fr)

	return &fr, result.Error
}

func (r *friendRepository) FindByID(id string) (*model.User, error) {
	user := &model.User{}

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.
		Preload("Friends").
		Preload("Requests").
		Where("id = ?", id).
		First(&user).Error;
		err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("uid", id)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

func (r *friendRepository) DeleteRequest(memberId string, userId string) error {
	return r.DB.Exec("DELETE FROM friend_requests WHERE receiver_id = ? AND sender_id = ?", memberId, userId).Error
}

func (r *friendRepository) RemoveFriend(memberId string, userId string) error {
	return r.DB.
		Exec("DELETE FROM friends WHERE user_id = ? AND friend_id = ?", memberId, userId).
		Exec("DELETE FROM friends WHERE user_id = ? AND friend_id = ?", userId, memberId).Error
}

func (r *friendRepository) Save(user *model.User) error {
	return r.DB.Save(&user).Error
}
