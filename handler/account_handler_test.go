package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMe(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		uid, _ := service.GenerateId()

		mockUserResp := &model.User{
			Email:    "bob@bob.com",
			Username: "Bobby",
			Image:    "image",
			IsOnline: false,
		}
		mockUserResp.ID = uid

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUserResp, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use a middleware to set context for test
		// the only claims we care about in this test
		// is the UID
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("userId", uid)
		})

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockUserResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t) // assert that UserService.Get was called
	})

	t.Run("NoContextUser", func(t *testing.T) {
		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", mock.Anything).Return(nil, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// do not append user to context
		router := gin.Default()
		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockUserService.AssertNotCalled(t, "Get", mock.Anything)
	})

	t.Run("NotFound", func(t *testing.T) {
		uid, _ := service.GenerateId()
		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(nil, fmt.Errorf("Some error down call chain"))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("userId", uid)
		})

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respErr := apperrors.NewNotFound("user", uid)

		respBody, err := json.Marshal(gin.H{
			"error": respErr,
		})
		assert.NoError(t, err)

		assert.Equal(t, respErr.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t) // assert that UserService.Get was called
	})
}

func TestEdit(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	uid, _ := service.GenerateId()
	user := &model.User{}
	user.ID = uid

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("userId", uid)
	})

	mockUserService := new(mocks.UserService)
	mockUserService.On("Get", uid).Return(user, nil)

	NewHandler(&Config{
		R:           router,
		UserService: mockUserService,
	})

	t.Run("Data binding error", func(t *testing.T) {
		rr := httptest.NewRecorder()

		form := url.Values{}
		form.Add("email", "notanemail")
		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "UpdateAccount")
	})

	t.Run("Update success", func(t *testing.T) {
		rr := httptest.NewRecorder()

		newName := "Sen"
		newEmail := "sen@example.com"

		form := url.Values{}
		form.Add("username", newName)
		form.Add("email", newEmail)

		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))
		request.Form = form

		userToUpdate := &model.User{
			Username: newName,
			Email:    newEmail,
		}
		userToUpdate.ID = user.ID

		updateArgs := mock.Arguments{
			userToUpdate,
		}

		dbImageURL := "https://jacobgoodwin.me/static/696292a38f493a4283d1a308e4a11732/84d81/Profile.jpg"

		mockUserService.
			On("IsEmailAlreadyInUse", newEmail).Return(false)

		mockUserService.
			On("UpdateAccount", updateArgs...).
			Run(func(args mock.Arguments) {
				userArg := args.Get(0).(*model.User)
				userArg.Image = dbImageURL
			}).
			Return(nil)

		router.ServeHTTP(rr, request)

		userToUpdate.Image = dbImageURL
		respBody, _ := json.Marshal(userToUpdate)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "IsEmailAlreadyInUse", newEmail)
		mockUserService.AssertCalled(t, "UpdateAccount", updateArgs...)
	})

	t.Run("Update failure", func(t *testing.T) {
		rr := httptest.NewRecorder()

		uid, _ := service.GenerateId()
		user := &model.User{}
		user.ID = uid

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("userId", uid)
		})

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(user, nil)

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		newName := "Sen"
		newEmail := "sen@example.com"

		form := url.Values{}
		form.Add("username", newName)
		form.Add("email", newEmail)

		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))
		request.Form = form

		userToUpdate := &model.User{
			Username: newName,
			Email:    newEmail,
		}
		userToUpdate.ID = user.ID

		updateArgs := mock.Arguments{
			userToUpdate,
		}

		mockError := apperrors.NewInternal()

		mockUserService.
			On("IsEmailAlreadyInUse", newEmail).Return(false)

		mockUserService.
			On("UpdateAccount", updateArgs...).
			Return(mockError)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "IsEmailAlreadyInUse", newEmail)
		mockUserService.AssertCalled(t, "UpdateAccount", updateArgs...)
	})

	t.Run("Email already in use", func(t *testing.T) {
		rr := httptest.NewRecorder()

		newName := "Sen"
		newEmail := "duplicate@example.com"

		form := url.Values{}
		form.Add("username", newName)
		form.Add("email", newEmail)

		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))
		request.Form = form

		userToUpdate := &model.User{
			Username: newName,
			Email:    newEmail,
		}
		userToUpdate.ID = user.ID

		updateArgs := mock.Arguments{
			userToUpdate,
		}

		mockUserService.
			On("IsEmailAlreadyInUse", newEmail).Return(true)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"field":   "Email",
			"message": "email already in use",
		})

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "IsEmailAlreadyInUse", newEmail)
		mockUserService.AssertNotCalled(t, "UpdateAccount", updateArgs...)
	})

	t.Run("Username too short", func(t *testing.T) {
		rr := httptest.NewRecorder()

		tooShortName := "Se"
		newEmail := "sen@example.com"

		form := url.Values{}
		form.Add("username", tooShortName)
		form.Add("email", newEmail)

		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))
		request.Form = form

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "UpdateAccount")
	})

	t.Run("Username too long", func(t *testing.T) {
		rr := httptest.NewRecorder()

		tooShortName := "Seoiasoidhaoushgduasgdiuagsdziuagszidgas"
		newEmail := "sen@example.com"

		form := url.Values{}
		form.Add("username", tooShortName)
		form.Add("email", newEmail)

		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))
		request.Form = form

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "UpdateAccount")
	})
}
