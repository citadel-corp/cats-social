package cat

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/citadel-corp/cats-social/internal/common/id"
)

var (
	AgeRegex = regexp.MustCompile("[<>]*\\d+")
)

type Service interface {
	List(ctx context.Context, req ListCatPayload, userID int) ([]CatResponse, error)
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

// List implements Service.
func (s *userService) List(ctx context.Context, req ListCatPayload, userID int) ([]CatResponse, error) {
	// validate payload; if invalid, set to empty so it will be ignored when querying
	if req.Limit == 0 {
		req.Limit = 5
	}
	if !slices.Contains(CatRaces, CatRace(req.Race)) {
		req.Race = ""
	}
	if !slices.Contains(CatSexes, CatSex(req.Sex)) {
		req.Sex = ""
	}
	req.AgeSearchType = IgnoreAge
	ageStr := AgeRegex.FindString(req.AgeInMonth)
	if ageStr != "" {
		searchType := rune(ageStr[0])
		switch searchType {
		case '>':
			req.AgeSearchType = MoreThan
			age, _ := strconv.Atoi(ageStr[1:])
			req.Age = age
		case '<':
			req.AgeSearchType = LessThan
			age, _ := strconv.Atoi(ageStr[1:])
			req.Age = age
		}
		if age, err := strconv.Atoi(ageStr); err == nil {
			req.Age = age
			req.AgeSearchType = EqualTo
		}
	}
	req.HasMatchedType = IgnoreHasMatched
	if req.HasMatched == "true" {
		req.HasMatchedType = HasMatched
	} else if req.HasMatched == "false" {
		req.HasMatchedType = HasNotMatched
	}
	cats, err := s.repository.List(ctx, req, userID)
	if err != nil {
		return nil, err
	}
	res := make([]CatResponse, len(cats))
	for i, cat := range cats {
		res[i] = CatResponse{
			ID:          cat.UID,
			Name:        cat.Name,
			Race:        string(cat.Race),
			Sex:         string(cat.Sex),
			AgeInMonth:  cat.Age,
			ImageUrls:   cat.ImageURLS,
			Description: cat.Description,
			HasMatched:  cat.HasMatched,
			CreatedAt:   cat.CreatedAt,
		}
	}
	return res, nil
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
