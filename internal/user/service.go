package user

import (
	"context"
	"fmt"
	"time"

	"github.com/citadel-corp/cats-social/internal/common/id"
	"github.com/citadel-corp/cats-social/internal/common/jwt"
	"github.com/citadel-corp/cats-social/internal/common/password"
)

type Service interface {
	Create(ctx context.Context, req CreateUserPayload) (*UserResponse, error)
	Login(ctx context.Context, req LoginPayload) (*UserResponse, error)
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s *userService) Create(ctx context.Context, req CreateUserPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	user := &User{
		UID:            id.GenerateStringID(16),
		Email:          req.Email,
		Name:           req.Name,
		HashedPassword: hashedPassword,
	}
	user, err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*8, fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		Email:       req.Email,
		Name:        req.Name,
		AccessToken: accessToken,
	}, nil
}

func (s *userService) Login(ctx context.Context, req LoginPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	match, err := password.Matches(req.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrWrongPassword
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*8, fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		Email:       user.Email,
		Name:        user.Name,
		AccessToken: accessToken,
	}, nil
}
