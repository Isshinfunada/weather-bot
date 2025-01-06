package usecase

import (
	"context"
	"fmt"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
)

type UserUsecase interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
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
