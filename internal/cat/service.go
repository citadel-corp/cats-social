package cat

import (
	"context"
	"fmt"

	"github.com/citadel-corp/cats-social/internal/common/id"
)

type Service interface {
	Create(ctx context.Context, req CreateCatPayload, userID string) (*CreateCatResponse, error)
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s *userService) Create(ctx context.Context, req CreateCatPayload, userID string) (*CreateCatResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	cat := &Cat{
		ID:          id.GenerateStringID(16),
		UserID:      userID,
		Name:        req.Name,
		Race:        CatRace(req.Race),
		Sex:         CatSex(req.Sex),
		Age:         req.AgeInMonth,
		Description: req.Description,
		HasMatched:  false,
		ImageURLS:   req.ImageURLS,
	}
	err = s.repository.Create(ctx, cat)
	if err != nil {
		return nil, err
	}
	return &CreateCatResponse{
		Id:        cat.ID,
		CreatedAt: cat.CreatedAt,
	}, nil
}
