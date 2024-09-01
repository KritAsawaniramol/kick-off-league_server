package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kickoff-league.com/authentication"
	"kickoff-league.com/handlers"
	model "kickoff-league.com/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"kickoff-league.com/usecases/authUsecase"
)

var secretKey = "JwtSecretKey"

func setupAuthRouter(httpHandler handlers.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	auth := authentication.NewJwtAuthentication(secretKey)

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register/organizer", httpHandler.RegisterOrganizer)
		authRouter.POST("/login", httpHandler.LoginUser)
		authRouter.POST("/register/normal", httpHandler.RegisterNormaluser)
		authRouter.POST("/logout", auth.Auth(), httpHandler.LogoutUser)
	}
	return router
}

func TestRegisterOrganizer(t *testing.T) {
	mockAuthUsecase := new(authUsecase.AuthUsecaseMock)
	httpHandler := handlers.NewhttpHandler(nil, nil, mockAuthUsecase, nil, nil, nil, nil, nil)
	router := setupAuthRouter(httpHandler)
	tests := []struct {
		name           string
		input          interface{}
		mockReturn     error
		expectedStatus int
	}{
		{
			name: "success",
			input: model.RegisterOrganizer{
				RegisterUser: model.RegisterUser{
					Email:    "test@example.com",
					Password: "password",
				},
				OrganizerName: "Test Organizer",
				Phone:         "1234567890",
			},
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "bad request - missing fields",
			input: model.RegisterOrganizer{
				RegisterUser: model.RegisterUser{
					Email: "test@example.com",
					// Password is missing
				},
				// OrganizerName and Phone are missing
			},
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failure - usecase error",
			input: model.RegisterOrganizer{
				RegisterUser: model.RegisterUser{
					Email:    "test@example.com",
					Password: "password",
				},
				OrganizerName: "Test Organizer",
				Phone:         "1234567890",
			},
			mockReturn:     errors.New("usecase error"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUsecase.ExpectedCalls = nil
			var body []byte
			if s, ok := tt.input.(string); ok {
				body = []byte(s) // use input as string directly for invalid JSON case
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req, _ := http.NewRequest(http.MethodPost, "/auth/register/organizer", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			mockAuthUsecase.On("RegisterOrganizer", mock.AnythingOfType("*model.RegisterOrganizer")).Return(tt.mockReturn)
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockAuthUsecase.AssertExpectations(t)
		})
	}
}

func TestRegisterNormaluser(t *testing.T) {
	mockAuthUsecase := new(authUsecase.AuthUsecaseMock)
	httpHandler := handlers.NewhttpHandler(nil, nil, mockAuthUsecase, nil, nil, nil, nil, nil)
	router := setupAuthRouter(httpHandler)

	tests := []struct {
		name           string
		input          interface{}
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful registration",
			input: model.RegisterNormaluser{
				RegisterUser: model.RegisterUser{
					Email:    "user@example.com",
					Password: "password123",
				},
				Username: "testuser",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "register success",
		},
		{
			name:           "bad request - JSON binding error",
			input:          `{"email": "user@example.com"}`, // missing password and username
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "required", // Check for required field error message
		},
		{
			name: "bad request - registration failure",
			input: model.RegisterNormaluser{
				RegisterUser: model.RegisterUser{
					Email:    "user@example.com",
					Password: "password123",
				},
				Username: "testuser",
			},
			mockError:      errors.New("registration failed"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "registration failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUsecase.ExpectedCalls = nil
			var body []byte
			if s, ok := tt.input.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req, _ := http.NewRequest(http.MethodPost, "/auth/register/normal", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			mockAuthUsecase.On("RegisterNormaluser", mock.AnythingOfType("*model.RegisterNormaluser")).Return(tt.mockError)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Contains(t, resp.Body.String(), tt.expectedBody)

			mockAuthUsecase.AssertExpectations(t)
		})
	}
}

func TestLoginUser(t *testing.T) {
	mockAuthUsecase := new(authUsecase.AuthUsecaseMock)
	httpHandler := handlers.NewhttpHandler(nil, nil, mockAuthUsecase, nil, nil, nil, nil, nil)
	router := setupAuthRouter(httpHandler)

	tests := []struct {
		name           string
		input          interface{}
		mockJWT        string
		mockUser       model.LoginResponse
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful login",
			input: model.LoginUser{
				Email:    "test@example.com",
				Password: "password",
			},
			mockJWT:        "mockToken",
			mockUser:       model.LoginResponse{ID: 1, Email: "test@example.com", Role: "organizer"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Login success",
		},
		{
			name:           "bad request - JSON binding error",
			input:          `{"email": "test@example.com"}`, // missing password
			mockJWT:        "",
			mockUser:       model.LoginResponse{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "required",
		},
		{
			name: "unauthorized - login failure",
			input: model.LoginUser{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockJWT:        "",
			mockUser:       model.LoginResponse{},
			mockError:      errors.New("invalid credentials"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUsecase.ExpectedCalls = nil
			var body []byte
			if s, ok := tt.input.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			mockAuthUsecase.On("Login", mock.AnythingOfType("*model.LoginUser")).Return(tt.mockJWT, tt.mockUser, tt.mockError)

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Contains(t, resp.Body.String(), tt.expectedBody)

			mockAuthUsecase.AssertExpectations(t)
		})
	}
}

func TestLogoutUser(t *testing.T) {

	mockAuthUsecase := new(authUsecase.AuthUsecaseMock)
	httpHandler := handlers.NewhttpHandler(nil, nil, mockAuthUsecase, nil, nil, nil, nil, nil)
	router := setupAuthRouter(httpHandler)

	// Create a valid token for testing
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 123, // Example user ID
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secretKey))

	req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)

	req.Header.Add("Cookie", fmt.Sprintf("token=%s", tokenString))

	resp := httptest.NewRecorder()

	// Execute the request
	router.ServeHTTP(resp, req)

	// Validate the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "token=; Path=/; Domain=localhost; Max-Age=0; HttpOnly", resp.Header().Get("Set-Cookie")) // Ensure token is cleared
	assert.Equal(t, "/home", resp.Header().Get("Location"))                                                   // Check redirection to /home
}
