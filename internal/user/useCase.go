package user

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"pikpo2/helpers/exception"
	"pikpo2/helpers/response"
	"pikpo2/internal/tier"
	"pikpo2/models"
)

type (
	UserUseCase interface {
		Register(ctx context.Context, user models.User) response.Response
		ReadOne(ctx context.Context, params int64) response.Response
		Update(ctx context.Context, ID int64, userID int64, user models.User) response.Response
	}

	userUseCaseImpl struct {
		repository UserRepository

		tierRepository tier.TierRepository
	}
)

func NewUserUseCase(repo UserRepository, tierRepo tier.TierRepository) UserUseCase {
	return &userUseCaseImpl{
		repository:     repo,
		tierRepository: tierRepo,
	}
}

func (nu *userUseCaseImpl) Register(ctx context.Context, user models.User) response.Response {

	var userResponse models.User
	var tier models.Tier

	min := 1
	max := 999
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(max-min) + min

	userResponse.Status = "free"
	userResponse.UsernameT3 = user.UsernameT3 + "." + strconv.Itoa(v)
	userResponse.Role = user.Role
	user.ExpiredAt = time.Now()

	ID, err := nu.repository.Create(ctx, userResponse)
	if err != nil {
		return response.Error(response.StatusInternalServerError, err)
	}

	userResponse.ID = int(ID)

	tier.Tier = "Tier-3"
	tier.UserId = userResponse.ID

	IDtier, err := nu.tierRepository.Create(ctx, tier)
	if err != nil {
		return response.Error(response.StatusInternalServerError, err)
	}

	tier.ID = int(IDtier)
	userResponse.Tier = append(userResponse.Tier, tier)

	return response.Success(response.StatusCreated, userResponse)
}

func (nu *userUseCaseImpl) ReadOne(ctx context.Context, params int64) response.Response {

	account, err := nu.repository.FindByID(ctx, params)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if err != nil {
		return response.Error(response.StatusInternalServerError, exception.ErrInternalServer)
	}

	return response.Success(response.StatusOK, account)
}

func (nu *userUseCaseImpl) Update(ctx context.Context, ID int64, userID int64, user models.User) response.Response {
	// var userResponse models.User
	// var tier models.Tier

	account, err := nu.repository.FindByID(ctx, ID)
	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if account.Status != "paid" {
		return response.Error(exception.ErrNotPremium.Error(), exception.ErrNotPremium)
	}

	accountTier, err := nu.tierRepository.FindByUserID(ctx, userID)

	if err == exception.ErrNotFound {
		return response.Error(response.StatusNotFound, exception.ErrNotFound)
	}

	if accountTier.Tier == "Tier-2" {
		account.UsernameT2 = user.UsernameT3 + "." + account.Role

		nu.repository.Update(ctx, ID, account)
	}

	if accountTier.Tier == "Tier-1" {
		account.UsernameT1 = user.UsernameT3

		nu.repository.Update(ctx, ID, account)
	}

	return response.Success(response.StatusOK, account)
}
