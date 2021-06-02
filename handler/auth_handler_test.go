package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Email, Username and Password Required", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "",
			"username": "",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Register")
	})

	t.Run("Invalid email", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob",
			"username": "bobby",
			"password": "supersecret1234",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Username too short", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"username": "bo",
			"password": "superpassword",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Register")
	})

	t.Run("Password too short", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"username": "bobby",
			"password": "supe",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Register")
	})

	t.Run("Username too long", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"username": "kjhasiudaiusdiuadiuagszuidgaiszugdziasgdiazgsdiazugdipas",
			"password": "superpassword",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Register")
	})

	t.Run("Error returned from UserService", func(t *testing.T) {
		u := &model.User{
			Email:    "bob@bob.com",
			Username: "bobby",
			Password: "avalidpassword",
		}

		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", u).Return(apperrors.NewConflict("User Already Exists", u.Email))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"username": u.Username,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 409, rr.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Successful Creation", func(t *testing.T) {
		u := &model.User{
			Email:    "bob@bob.com",
			Username: "bobby",
			Password: "avalidpassword",
		}

		mockUserService := new(mocks.UserService)

		mockUserService.
			On("Register", u).
			Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"username": u.Username,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(u)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// setup mock services, gin engine/router, handler layer
	mockUserService := new(mocks.UserService)

	router := gin.Default()

	NewHandler(&Config{
		R:           router,
		UserService: mockUserService,
	})

	t.Run("Bad request data", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with invalid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    "notanemail",
			"password": "short",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Login")
	})

	t.Run("Error Returned from UserService.Login", func(t *testing.T) {
		email := "bob@bob.com"
		password := "pwdoesnotmatch123"

		mockUSArgs := mock.Arguments{
			&model.User{Email: email, Password: password},
		}

		// so we can check for a known status code
		mockError := apperrors.NewAuthorization("invalid email/password combo")

		mockUserService.On("Login", mockUSArgs...).Return(mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockUserService.AssertCalled(t, "Login", mockUSArgs...)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Successful Login", func(t *testing.T) {
		email := "bob@bob.com"
		password := "pwworksgreat123"

		u := &model.User{
			Email:    email,
			Password: password,
		}

		mockUSArgs := mock.Arguments{u}

		mockUserService.On("Login", mockUSArgs...).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(u)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertCalled(t, "Login", mockUSArgs...)
	})

}
