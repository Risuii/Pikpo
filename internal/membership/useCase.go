package membership

import (
	"context"

	"pikpo2/helpers/exception"
	"pikpo2/helpers/response"
	"pikpo2/internal/tier"
	"pikpo2/internal/user"
	"pikpo2/models"
)

type (
	MembershipUseCase interface {
		Membership(ctx context.Context, ID int64, tierNumber int64) response.Response
	}

	membershipUseCaseImpl struct {
		repository     user.UserRepository
		tierRepository tier.TierRepository
	}
)

func NewMembershipUseCase(repo user.UserRepository, tierRepo tier.TierRepository) MembershipUseCase {
	return &membershipUseCaseImpl{
		repository: repo,

		tierRepository: tierRepo,
	}
}

func (nt *membershipUseCaseImpl) Membership(ctx context.Context, ID int64, tierNumber int64) response.Response {

	var tier models.Tier

	account, err := nt.repository.FindByID(ctx, ID)

	if err == exception.ErrNotFound {
		return response.Error(exception.ErrNotFound.Error(), exception.ErrBadRequest)
	}

	account.Status = "paid"

	_ = nt.repository.UpdateStatus(ctx, ID, account)

	_ = nt.tierRepository.UpdateTier(ctx, ID, tierNumber, tier)

	return nil
}
