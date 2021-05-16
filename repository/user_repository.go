package repository

import (
	"errors"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"gorm.io/gorm"
	"log"
	"regexp"
)

// userRepository is data/repository implementation
// of service layer UserRepository
type userRepository struct {
	DB *gorm.DB
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	user := &model.User{}

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("uid", id)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

func (r *userRepository) Create(u *model.User) error {
	if result := r.DB.Create(&u); result.Error != nil {
		// check unique constraint
		if isDuplicateKeyError(result.Error) {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, result.Error)
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, result.Error)
		return apperrors.NewInternal()
	}
	return nil
}

// FindByEmail retrieves user row by email address
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("email", email)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

func (r *userRepository) Update(u *model.User) error {
	result := r.DB.Save(&u)
	return result.Error
}

func isDuplicateKeyError(err error) bool {
	duplicate := regexp.MustCompile(`\(SQLSTATE 23505\)$`)
	return duplicate.MatchString(err.Error())
}
