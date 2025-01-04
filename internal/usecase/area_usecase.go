package usecase

import (
	"context"
	"fmt"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
)

type AreaUseCase interface {
	GetHierarchy(ctx context.Context, class20ID int) (*entity.HierarchyArea, error)
}

type areaUseCase struct {
	areaRepo repository.AreaRepository
}

func NewAreaUseCase(aRepo repository.AreaRepository) AreaUseCase {
	return &areaUseCase{areaRepo: aRepo}
}

func (u *areaUseCase) GetHierarchy(ctx context.Context, class20ID int) (*entity.HierarchyArea, error) {
	// バリデーション
	if len(fmt.Sprintf("%d", class20ID)) != 7 {
		return nil, fmt.Errorf("id length is invalid")
	}

	hierarchy, err := u.areaRepo.FindHierarchyByClass20ID(ctx, class20ID)
	if err != nil {
		return nil, err
	}
	if hierarchy == nil {
		return nil, fmt.Errorf("not found for class20 id=%d", class20ID)
	}

	// ここに通知ロジックや他レポジトリとの連携処理を追記できる
	return hierarchy, nil
}
