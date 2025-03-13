package service_test

import (
	"bookstore-framework/internal/users"
	"bookstore-framework/internal/users/api/dto"
	mocks "bookstore-framework/test/mock"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Success(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	jwtGen := mocks.NewMockJWTGenerator(ctrl)
	service := users.NewUserService(mockRepo, jwtGen)

	t.Run("Register", func(t *testing.T) {
		ctx := context.Background()
		req := dto.RegisterRequest{
			Username: "test",
			Name:     "testuser",
			Email:    "test@gmail.com",
			Password: "password123",
		}

		expectedUser := &users.User{
			ID:       1,
			Name:     req.Name,
			Username: req.Username,
			Email:    req.Email,
		}

		mockRepo.EXPECT().Register(gomock.Any(), gomock.Any()).Return(expectedUser, nil)

		result, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.Equal(t, req.Username, result.Username)
	})

	t.Run("Login", func(t *testing.T) {
		ctx := context.Background()
		password := "password"
		req := dto.LoginRequest{
			Username: "test",
			Password: password,
		}

		expectedToken := "mocked-jwt-token"
		expectedRes := &dto.LoginResponse{
			TokenAccess: expectedToken,
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		mockUser := &users.User{
			ID:       1,
			Username: req.Username,
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}
		mockRepo.EXPECT().FindUserByUsername(gomock.Any(), req.Username).Return(mockUser, nil)
		jwtGen.EXPECT().GenerateToken(mockUser.ID, mockUser.Username, mockUser.Email).Return(expectedToken, nil)

		result, err := service.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedRes.TokenAccess, result.TokenAccess)

	})

	t.Run("GetProfile", func(t *testing.T) {
		ctx := context.Background()

		mockUser := &users.User{
			ID:         1,
			Username:   "XXXX",
			Name:       "testuser",
			Email:      "XXXXXXXXXXXXXX",
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}

		mockRepo.EXPECT().FindUserByID(gomock.Any(), uint(1)).Return(mockUser, nil)

		result, err := service.GetProfile(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockUser.ID, result.ID)
		assert.Equal(t, mockUser.Username, result.Username)

	})

}

func TestUserService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	jwtGen := mocks.NewMockJWTGenerator(ctrl)
	service := users.NewUserService(mockRepo, jwtGen)

	t.Run("Register", func(t *testing.T) {
		ctx := context.Background()
		req := dto.RegisterRequest{
			Username: "test",
			Name:     "testuser",
			Email:    "test@gmail.com",
			Password: "password123",
		}
		mockRepo.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error"))

		result, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "Error", err.Error())
	})

	t.Run("Login", func(t *testing.T) {
		ctx := context.Background()

		req := dto.LoginRequest{
			Username: "test",
			Password: "wrongpassword",
		}
		correctPassword := "correctpassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
		assert.NoError(t, err)

		mockUser := &users.User{
			ID:       1,
			Username: req.Username,
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}
		mockRepo.EXPECT().FindUserByUsername(gomock.Any(), req.Username).Return(mockUser, nil)

		result, err := service.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid password or username", err.Error())

	})

	t.Run("GetProfile", func(t *testing.T) {
		ctx := context.Background()

		mockRepo.EXPECT().FindUserByID(gomock.Any(), uint(1)).Return(nil, errors.New("User not found"))

		result, err := service.GetProfile(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "User not found", err.Error())

	})
}
