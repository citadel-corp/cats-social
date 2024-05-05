package catmatch

import (
	"context"

	"github.com/citadel-corp/cats-social/internal/cat"
	"github.com/citadel-corp/cats-social/internal/common/id"
)

type Service interface {
	Create(ctx context.Context, req PostCatMatch, userID int64) error
	Approve(ctx context.Context, req ApproveOrRejectMatch, userId int64) error
	Reject(ctx context.Context, req ApproveOrRejectMatch, userId int64) error
	Delete(ctx context.Context, req DeleteMatch, userId int64) error
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
	issuerCat, err := s.catRepository.GetByUIDAndUserID(ctx, req.UserCatId, userID)
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

func (s *catMatchService) Approve(ctx context.Context, req ApproveOrRejectMatch, userId int64) error {
	// get match
	filter := map[string]interface{}{
		"pending_only": true,
	}
	match, err := s.repository.GetByUIDAndUserID(ctx, req.MatchUID, userId, filter)
	if err != nil {
		return err
	}

	err = s.repository.Approve(ctx, match)
	if err != nil {
		return err
	}

	return nil
}

// List implements Service.
func (s *catMatchService) List(ctx context.Context, userID int64) ([]CatMatchResponse, error) {
	filter := map[string]interface{}{}
	catMatches, err := s.repository.List(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	return MakeCatMatchResponse(catMatches, userID), nil
}

func (s *catMatchService) Reject(ctx context.Context, req ApproveOrRejectMatch, userId int64) error {
	// get match
	filter := map[string]interface{}{
		"pending_only": true,
	}
	match, err := s.repository.GetByUIDAndUserID(ctx, req.MatchUID, userId, filter)
	if err != nil {
		return err
	}

	err = s.repository.Reject(ctx, match)
	if err != nil {
		return err
	}

	return nil
}

func (s *catMatchService) Delete(ctx context.Context, req DeleteMatch, userId int64) error {
	// get match
	filter := map[string]interface{}{
		"pending_only": true,
	}
	match, err := s.repository.GetByUIDAndUserID(ctx, req.MatchUID, userId, filter)
	if err != nil {
		return err
	}

	err = s.repository.Delete(ctx, match.ID, userId)
	if err != nil {
		return err
	}

	return nil
}
