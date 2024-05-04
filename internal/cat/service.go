package cat

import (
	"context"
	"fmt"
	"strings"

	"github.com/citadel-corp/cats-social/internal/common/id"
)

type Service interface {
	Create(ctx context.Context, req CreateUpdateCatPayload, userID int) (*CreateCatResponse, error)
	Update(ctx context.Context, req CreateUpdateCatPayload, id string, userID int) error
	Delete(ctx context.Context, id string, userID int) error
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s *userService) Create(ctx context.Context, req CreateUpdateCatPayload, userID int) (*CreateCatResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	cat := &Cat{
		UID:         id.GenerateStringID(16),
		UserID:      userID,
		Name:        req.Name,
		Race:        CatRace(req.Race),
		Sex:         CatSex(req.Sex),
		Age:         req.AgeInMonth,
		Description: req.Description,
		HasMatched:  false,
		ImageURLS:   req.ImageURLS,
	}
	cat, err = s.repository.Create(ctx, cat)
	if err != nil {
		return nil, err
	}
	return &CreateCatResponse{
		Id:        cat.UID,
		CreatedAt: cat.CreatedAt,
	}, nil
}

// Update implements Service.
func (s *userService) Update(ctx context.Context, req CreateUpdateCatPayload, uid string, userID int) error {
	err := req.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	cat, err := s.repository.GetByUIDAndUserID(ctx, uid, userID)
	if err != nil {
		return err
	}

	if cat.Sex != req.Sex && cat.HasMatched {
		return ErrCatHasMatched
	}

	cat = &Cat{
		UID:         uid,
		UserID:      userID,
		Name:        req.Name,
		Race:        req.Race,
		Sex:         CatSex(strings.ToLower(string(req.Sex))),
		Age:         req.AgeInMonth,
		Description: req.Description,
		HasMatched:  cat.HasMatched,
		ImageURLS:   req.ImageURLS,
	}
	return s.repository.Update(ctx, cat)

}

// Delete implements Service.
func (s *userService) Delete(ctx context.Context, id string, userID int) error {
	return s.repository.Delete(ctx, id, userID)
}
