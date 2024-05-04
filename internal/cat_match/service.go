package catmatch

import (
	"context"

	"github.com/citadel-corp/cats-social/internal/cat"
	"github.com/citadel-corp/cats-social/internal/common/id"
)

type Service interface {
	Create(ctx context.Context, req PostCatMatch, userID int64) error
	List(ctx context.Context, userID int64) ([]CatMatchResponse, error)
}

type catMatchService struct {
	repository    Repository
	catRepository cat.Repository
}

func NewService(repository Repository, catRepository cat.Repository) Service {
	return &catMatchService{repository: repository, catRepository: catRepository}
}

func (s *catMatchService) Create(ctx context.Context, req PostCatMatch, userID int64) error {
	err := req.Validate()
	if err != nil {
		return ErrValidationFailed
	}

	// get issuer cat
	issuerCat, err := s.catRepository.GetByUIDAndUserID(ctx, req.UserCatId, int(userID))
	if err != nil {
		return err
	}

	// get matched cat
	matchedCat, err := s.catRepository.GetByUID(ctx, req.MatchCatId)
	if err != nil {
		return err
	}

	if issuerCat.UserID == matchedCat.UserID {
		return ErrCatSameUser
	}

	if issuerCat.Sex == matchedCat.Sex {
		return ErrCatSameSex
	}

	if issuerCat.HasMatched || matchedCat.HasMatched {
		return ErrCatHasMatched
	}

	catMatch := &CatMatches{
		UID:         id.GenerateStringID(16),
		IssuerCatId: issuerCat.ID,
		IssueUserId: userID,
		MatchCatId:  matchedCat.ID,
		MatchUserId: int64(matchedCat.UserID),
		Message:     req.Message,
	}
	err = s.repository.Create(ctx, catMatch)
	if err != nil {
		return err
	}

	return nil
}

// List implements Service.
func (s *catMatchService) List(ctx context.Context, userID int64) ([]CatMatchResponse, error) {
	catMatches, err := s.repository.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]CatMatchResponse, len(catMatches))
	for i, catMatch := range catMatches {
		res[i] = CatMatchResponse{
			ID:             catMatch.UID,
			IssuedBy:       Issuer{},          // todo
			MatchCatDetail: cat.CatResponse{}, // todo
			UserCatDetail:  cat.CatResponse{}, // todo
			Message:        catMatch.Message,
			CreatedAt:      *catMatch.CreatedAt,
		}
	}
	return res, nil
}
