package handler_test

import (
	"bookstore-framework/internal/api"
	"bookstore-framework/internal/api/dto"
	mocks "bookstore-framework/test/mock"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fahrizalvianaz/shared-response/httputil"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := api.NewUserHandler(mockService)

	t.Run("Register", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username: "test",
			Name:     "testuser",
			Email:    "test@gmail.com",
			Password: "password123",
		}

		data := dto.RegisterResponse{
			ID:        1,
			Username:  req.Username,
			Email:     req.Email,
			CreatedAt: time.Now(),
		}

		mockService.EXPECT().Register(gomock.Any(), gomock.Eq(req)).
			Return(&data, nil)

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterHandler(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response httputil.Response
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, true, response.Status)
		assert.Equal(t, "User registered successfully", response.Message)

		var registerReponse dto.RegisterResponse

		dataBytes, _ := json.Marshal(response.Data)
		json.Unmarshal(dataBytes, &registerReponse)

		assert.Equal(t, data.ID, registerReponse.ID)
		assert.Equal(t, data.Username, registerReponse.Username)
	})

	t.Run("Login", func(t *testing.T) {
		req := dto.LoginRequest{
			Username: "test",
			Password: "test123",
		}
		accessToken := "access_token_jwt"
		res := dto.LoginResponse{
			TokenAccess: accessToken,
		}

		mockService.EXPECT().Login(gomock.Any(), gomock.Eq(req)).
			Return(&res, nil)

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response httputil.Response
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, true, response.Status)
		assert.Equal(t, "Login Successfully", response.Message)

		var login dto.LoginResponse

		dataBytes, _ := json.Marshal(response.Data)
		json.Unmarshal(dataBytes, &login)

		assert.Equal(t, accessToken, login.TokenAccess)
	})

	t.Run("GetProfile", func(t *testing.T) {
		res := &dto.ProfileResponse{
			ID:         1,
			Name:       "test",
			Username:   "testuser",
			Email:      "test@gmail.com",
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}

		mockService.EXPECT().GetProfile(gomock.Any(), uint(1)).
			Return(res, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		c.Set("userID", uint(1))

		handler.GetProfile(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response httputil.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, true, response.Status)
		assert.Equal(t, "Profile retrieve successfully", response.Message)

		var profileResponse dto.ProfileResponse
		dataBytes, _ := json.Marshal(response.Data)
		json.Unmarshal(dataBytes, &profileResponse)

		assert.Equal(t, res.ID, profileResponse.ID)
		assert.Equal(t, res.Username, profileResponse.Username)
		assert.Equal(t, res.Name, profileResponse.Name)
		assert.Equal(t, res.Email, profileResponse.Email)

	})
}

func TestUserHandler_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := api.NewUserHandler(mockService)

	t.Run("Register_MissingRequiredField", func(t *testing.T) {
		reqMissingField := map[string]interface{}{
			"username": "test",
			"email":    "test@gmail.com",
			"password": "password123",
		}

		body, err := json.Marshal(reqMissingField)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response httputil.Response
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, false, response.Status)
		assert.Equal(t, "Invalid Request format", response.Message)

		errorMsg, ok := response.Data.(string)
		if ok {
			assert.Contains(t, errorMsg, "Name")
			assert.Contains(t, errorMsg, "required")
		}

	})

	t.Run("Register_ServiceError", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username: "XXXX",
			Name:     "testuser",
			Email:    "XXXXXXXXXXXXXX",
			Password: "XXXXXXXXXXX",
		}

		errorMsg := "Email already exist"
		mockService.EXPECT().Register(gomock.Any(), gomock.Eq(req)).
			Return(nil, errors.New(errorMsg))

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response httputil.Response
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, false, response.Status)
		assert.Equal(t, errorMsg, response.Message)

	})

	t.Run("Login_MissingRequiredField", func(t *testing.T) {
		reqMissingField := map[string]interface{}{
			// "username": "test",
			"password": "password123",
		}

		body, err := json.Marshal(reqMissingField)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RegisterHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response httputil.Response
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, false, response.Status)
		assert.Equal(t, "Invalid Request format", response.Message)

		errorMsg, ok := response.Data.(string)
		if ok {
			assert.Contains(t, errorMsg, "Name")
			assert.Contains(t, errorMsg, "required")
		}

	})

	t.Run("Login_ServiceError", func(t *testing.T) {
		req := dto.LoginRequest{
			Username: "test",
			Password: "test123",
		}

		errorMsg := "Invalid username or password"

		mockService.EXPECT().Login(gomock.Any(), gomock.Eq(req)).
			Return(nil, errors.New(errorMsg))

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.LoginHandler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response httputil.Response
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, false, response.Status)
		assert.Equal(t, errorMsg, response.Message)

	})

	t.Run("GetProfile_JWTError", func(t *testing.T) {

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		// c.Set("userID", uint(1))

		handler.GetProfile(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response httputil.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Equal(t, false, response.Status)
		assert.Equal(t, "User not found", response.Message)

	})

	t.Run("GetProfile_ServiceError", func(t *testing.T) {

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		c.Set("userID", uint(1))

		mockService.EXPECT().GetProfile(gomock.Any(), uint(1)).
			Return(nil, errors.New("Failed to retrieve profile"))

		handler.GetProfile(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response httputil.Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, false, response.Status)
		assert.Equal(t, "Failed to retrieve profile", response.Message)

	})

}
