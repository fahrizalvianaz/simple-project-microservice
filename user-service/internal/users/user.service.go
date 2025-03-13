package users

import (
	"bookstore-framework/internal/users/api/dto"
	"bookstore-framework/pkg"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userId uint) (*dto.ProfileResponse, error)
}

type userService struct {
	userRepo UserRepository
	jwtGen   pkg.JWTGenerator
}

func NewUserService(userRepo UserRepository, jwtGen pkg.JWTGenerator) UserService {
	return &userService{
		userRepo: userRepo,
		jwtGen:   jwtGen,
	}
}

func (s *userService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := User{
		Username: req.Username,
		Name:     req.Name,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	registerUser, err := s.userRepo.Register(ctx, &user)
	if err != nil {
		return nil, err
	}

	respone := &dto.RegisterResponse{
		ID:        registerUser.ID,
		Username:  registerUser.Username,
		Email:     registerUser.Email,
		CreatedAt: registerUser.CreatedAt,
	}

	return respone, nil

}

func (s *userService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password or username")
	}

	token, err := s.jwtGen.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}
	respose := &dto.LoginResponse{
		TokenAccess: token,
	}

	return respose, nil

}

func (s *userService) GetProfile(ctx context.Context, userId uint) (*dto.ProfileResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := &dto.ProfileResponse{
		ID:         user.ID,
		Username:   user.Username,
		Name:       user.Name,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt,
		ModifiedAt: user.ModifiedAt,
	}

	return response, nil
}
