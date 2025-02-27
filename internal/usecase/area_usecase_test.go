package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAreaRepo は AreaRepository インターフェースのモック
type MockAreaRepo struct {
	mock.Mock
}

func (m *MockAreaRepo) FindHierarchyByClass20ID(ctx context.Context, class20ID string) (*entity.HierarchyArea, error) {
	args := m.Called(ctx, class20ID)
	// nilポインタの可能性を考慮
	var hier *entity.HierarchyArea
	if val := args.Get(0); val != nil {
		hier = val.(*entity.HierarchyArea)
	}
	return hier, args.Error(1)
}

// テスト用セットアップ関数
func setupAreaUsecaseTest() (*MockAreaRepo, usecase.AreaUseCase) {
	mockRepo := new(MockAreaRepo)
	areaUC := usecase.NewAreaUseCase(mockRepo)
	return mockRepo, areaUC
}

// 無効なID長のテスト
func TestAreaUsecase_GetHierarchy_InvalidIDLength(t *testing.T) {
	_, areaUC := setupAreaUsecaseTest()
	ctx := context.Background()

	// 長さが7でないIDを指定（例えば、12345は5桁）
	hierarchy, err := areaUC.GetHierarchy(ctx, "12345")
	assert.Nil(t, hierarchy)
	assert.EqualError(t, err, "id length is invalid")
}

// Hierarchyが見つからない場合のテスト
func TestAreaUsecase_GetHierarchy_NotFound(t *testing.T) {
	mockRepo, areaUC := setupAreaUsecaseTest()
	ctx := context.Background()

	validID := "1234567" // 長さ7の有効なID
	mockRepo.On("FindHierarchyByClass20ID", ctx, validID).Return(nil, nil)

	hierarchy, err := areaUC.GetHierarchy(ctx, validID)
	assert.Nil(t, hierarchy)
	assert.EqualError(t, err, fmt.Sprintf("not found for class20 id=%s", validID))
	mockRepo.AssertExpectations(t)
}

// リポジトリクエリエラー時のテスト
func TestAreaUsecase_GetHierarchy_RepoError(t *testing.T) {
	mockRepo, areaUC := setupAreaUsecaseTest()
	ctx := context.Background()

	validID := "1234567"
	repoErr := errors.New("database error")
	mockRepo.On("FindHierarchyByClass20ID", ctx, validID).Return(nil, repoErr)

	hierarchy, err := areaUC.GetHierarchy(ctx, validID)
	assert.Nil(t, hierarchy)
	assert.EqualError(t, err, "database error")
	mockRepo.AssertExpectations(t)
}

// 正常系のテスト
func TestAreaUsecase_GetHierarchy_Success(t *testing.T) {
	mockRepo, areaUC := setupAreaUsecaseTest()
	ctx := context.Background()

	validID := "1234567"
	expectedHierarchy := &entity.HierarchyArea{
		Class20: &entity.AreaClass20{ID: validID, Name: "Area20"},
		// 必要に応じて他の階層も設定
	}

	mockRepo.On("FindHierarchyByClass20ID", ctx, validID).Return(expectedHierarchy, nil)

	hierarchy, err := areaUC.GetHierarchy(ctx, validID)
	assert.NoError(t, err)
	assert.Equal(t, expectedHierarchy, hierarchy)
	mockRepo.AssertExpectations(t)
}
