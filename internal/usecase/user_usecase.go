package usecase

import (
	"context"
	"fmt"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
)

type UserUsecase interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, userID int) (*entity.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(ur repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: ur,
	}
}

// 新規ユーザー作成
func (u *userUsecase) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if user.LINEUserID == "" {
		return nil, fmt.Errorf("LINEUserID is required")
	}
	created, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// IDでユーザー検索
func (u *userUsecase) GetByID(ctx context.Context, userID int) (*entity.User, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	user, err := u.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found (id=%d)", userID)
	}
	return user, nil
}
