package usecase

import (
	"context"
	"fmt"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
)

type AreaUseCase interface {
	GetHierarchy(ctx context.Context, class20ID string) (*entity.HierarchyArea, error)
}

type areaUseCase struct {
	areaRepo repository.AreaRepository
}

func NewAreaUseCase(aRepo repository.AreaRepository) AreaUseCase {
	return &areaUseCase{areaRepo: aRepo}
}

func (u *areaUseCase) GetHierarchy(ctx context.Context, class20ID string) (*entity.HierarchyArea, error) {
	const CLASS20_LENGTH = 7

	if err := validateClass20ID(class20ID, CLASS20_LENGTH); err != nil {
		return nil, err
	}

	hierarchy, err := u.areaRepo.FindHierarchyByClass20ID(ctx, class20ID)
	if err != nil {
		return nil, err
	}
	if hierarchy == nil {
		return nil, fmt.Errorf("not found for class20 id=%s", class20ID)
	}

	return hierarchy, nil
}

func validateClass20ID(class20ID string, length int) error {
	if len(class20ID) != length {
		return fmt.Errorf("id length is invalid")
	}
	return nil
}
