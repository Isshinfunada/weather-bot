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
	GetByLINEID(ctx context.Context, LINEUserID string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
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

// LINEUserIDで探す
func (u *userUsecase) GetByLINEID(ctx context.Context, LINEUserID string) (*entity.User, error) {
	if LINEUserID == "" {
		return nil, fmt.Errorf("LINEUserID is required")
	}
	user, err := u.userRepo.FindUserByLINEUserID(ctx, LINEUserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found (LINEUserID=%s)", LINEUserID)
	}
	return user, nil
}

// 既存ユーザーの更新
func (u *userUsecase) Update(ctx context.Context, user *entity.User) error {
	if user.ID <= 0 {
		return fmt.Errorf("invalid user id")
	}
	// TODO: バリデーション
	err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
