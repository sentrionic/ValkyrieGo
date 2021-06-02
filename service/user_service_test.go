package service

import (
	"fmt"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/model/fixture"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUserResp := &model.User{
			Email:    "bob@bob.com",
			Username: "Bobby",
			Image:    "image",
			IsOnline: false,
		}
		mockUserResp.ID = uid

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", uid).Return(mockUserResp, nil)

		u, err := us.Get(uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResp)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", uid).Return(nil, fmt.Errorf("some error down the call chain"))

		u, err := us.Get(uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestRegister(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUser := &model.User{
			Email:    "bob@bob.com",
			Username: "bobby",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(0).(*model.User) // arg 0 is context, arg 1 is *User
				userArg.ID = uid
			}).Return(nil)

		err := us.Register(mockUser)

		assert.NoError(t, err)

		// assert user now has a userID
		assert.Equal(t, uid, mockUser.ID)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "bob@bob.com",
			Username: "bobby",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.UserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockErr := apperrors.NewConflict("email", mockUser.Email)

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mockUser).
			Return(mockErr)

		err := us.Register(mockUser)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())

		mockUserRepository.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	// setup valid email/pw combo with hashed password to test method
	// response when provided password is invalid
	email := "bob@bob.com"
	validPW := "howdyhoneighbor!"
	hashedValidPW, _ := hashPassword(validPW)
	invalidPW := "howdyhodufus!"

	mockUserRepository := new(mocks.UserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUser := &model.User{
			Email:    email,
			Password: validPW,
		}

		mockUserResp := &model.User{
			Email:    email,
			Password: hashedValidPW,
		}
		mockUserResp.ID = uid

		mockArgs := mock.Arguments{
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		err := us.Login(mockUser)

		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})

	t.Run("Invalid email/password combination", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUser := &model.User{
			Email:    email,
			Password: invalidPW,
		}

		mockUserResp := &model.User{
			Email:    email,
			Password: hashedValidPW,
		}
		mockUserResp.ID = uid

		mockArgs := mock.Arguments{
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		err := us.Login(mockUser)

		assert.Error(t, err)
		assert.EqualError(t, err, "Invalid email and password combination")
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})
}

func TestUpdateUser(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockFileRepository := new(mocks.FileRepository)

	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
		FileRepository: mockFileRepository,
	})

	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUser := &model.User{
			Email:    "new@bob.com",
			Username: "robert",
		}
		mockUser.ID = uid

		mockArgs := mock.Arguments{
			mockUser,
		}

		mockUserRepository.
			On("Update", mockArgs...).Return(nil)

		err := us.UpdateAccount(mockUser)

		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "Update", mockArgs...)
	})

	t.Run("Failure", func(t *testing.T) {
		uid, _ := GenerateId()

		mockUser := &model.User{}
		mockUser.ID = uid

		mockArgs := mock.Arguments{
			mockUser,
		}

		mockError := apperrors.NewInternal()

		mockUserRepository.
			On("Update", mockArgs...).Return(mockError)

		err := us.UpdateAccount(mockUser)
		assert.Error(t, err)

		apperror, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperrors.Internal, apperror.Type)

		mockUserRepository.AssertCalled(t, "Update", mockArgs...)
	})

	t.Run("Successful new image", func(t *testing.T) {
		uid, _ := GenerateId()

		// does not have have imageURL
		mockUser := &model.User{
			Email:    "new@bob.com",
			Username: "NewRobert",
		}
		mockUser.ID = uid

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "test_dir"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
		}

		imageURL := "https://imageurl.com/jdfkj34kljl"

		mockFileRepository.
			On("UploadAvatar", uploadFileArgs...).
			Return(imageURL, nil)

		updateArgs := mock.Arguments{
			mockUser,
		}

		mockUpdatedUser := &model.User{
			Email:    "new@bob.com",
			Username: "NewRobert",
			Image:    imageURL,
		}
		mockUpdatedUser.ID = uid

		mockUserRepository.
			On("Update", updateArgs...).
			Return(nil)

		url, err := us.ChangeAvatar(imageFileHeader, directory)
		assert.NoError(t, err)
		mockUser.Image = url

		err = us.UpdateAccount(mockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUpdatedUser, mockUser)
		mockFileRepository.AssertCalled(t, "UploadAvatar", uploadFileArgs...)
		mockUserRepository.AssertCalled(t, "Update", updateArgs...)
	})

	t.Run("Successful update image", func(t *testing.T) {
		imageURL := "https://imageurl.com/jdfkj34kljl"
		uid, _ := GenerateId()

		mockUser := &model.User{
			Email:    "new@bob.com",
			Username: "NewRobert",
			Image:    imageURL,
		}
		mockUser.ID = uid

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "test_dir"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
		}

		deleteImageArgs := mock.Arguments{
			imageURL,
		}

		mockFileRepository.
			On("UploadAvatar", uploadFileArgs...).
			Return(imageURL, nil)

		mockFileRepository.
			On("DeleteImage", deleteImageArgs...).
			Return(nil)

		mockUpdatedUser := &model.User{
			Email:    "new@bob.com",
			Username: "NewRobert",
			Image:    imageURL,
		}
		mockUpdatedUser.ID = uid

		updateArgs := mock.Arguments{
			mockUser,
		}

		mockUserRepository.
			On("Update", updateArgs...).
			Return(nil)

		url, err := us.ChangeAvatar(imageFileHeader, directory)
		assert.NoError(t, err)
		err = us.DeleteImage(mockUser.Image)
		assert.NoError(t, err)

		mockUser.Image = url
		err = us.UpdateAccount(mockUser)
		assert.NoError(t, err)

		assert.Equal(t, mockUpdatedUser, mockUser)
		mockFileRepository.AssertCalled(t, "UploadAvatar", uploadFileArgs...)
		mockFileRepository.AssertCalled(t, "DeleteImage", imageURL)
		mockUserRepository.AssertCalled(t, "Update", updateArgs...)
	})

	t.Run("FileRepository Error", func(t *testing.T) {
		// need to create a new UserService and repository
		// because testify has no way to overwrite a mock's
		// "On" call.
		mockUserRepository := new(mocks.UserRepository)
		mockFileRepository := new(mocks.FileRepository)

		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
			FileRepository: mockFileRepository,
		})

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "file_directory"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
		}

		mockError := apperrors.NewInternal()
		mockFileRepository.
			On("UploadAvatar", uploadFileArgs...).
			Return("", mockError)

		url, err := us.ChangeAvatar(imageFileHeader, directory)
		assert.Equal(t, "", url)
		assert.Error(t, err)

		mockFileRepository.AssertCalled(t, "UploadAvatar", uploadFileArgs...)
		mockUserRepository.AssertNotCalled(t, "Update")
	})

	t.Run("UserRepository UpdateImage Error", func(t *testing.T) {
		uid, _ := GenerateId()
		imageURL := "https://imageurl.com/jdfkj34kljl"

		// has imageURL
		mockUser := &model.User{
			Email:    "new@bob.com",
			Username: "A New Bob!",
			Image:    imageURL,
		}
		mockUser.ID = uid

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := "file_dir"

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
		}

		mockFileRepository.
			On("UploadAvatar", uploadFileArgs...).
			Return(imageURL, nil)

		updateArgs := mock.Arguments{
			mockUser,
		}

		mockError := apperrors.NewInternal()
		mockUserRepository.
			On("Update", updateArgs...).
			Return(mockError)

		url, err := us.ChangeAvatar(imageFileHeader, directory)
		assert.NoError(t, err)
		assert.Equal(t, imageURL, url)

		err = us.UpdateAccount(mockUser)

		assert.Error(t, err)
		mockFileRepository.AssertCalled(t, "UploadAvatar", uploadFileArgs...)
		mockUserRepository.AssertCalled(t, "Update", updateArgs...)
	})

}
