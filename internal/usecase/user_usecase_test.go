package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/Isshinfunada/weather-bot/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モックリポジトリ定義
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) FindUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if u := args.Get(0); u != nil {
		return u.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) FindUserByLINEUserID(ctx context.Context, LINEUserID string) (*entity.User, error) {
	args := m.Called(ctx, LINEUserID)
	if u := args.Get(0); u != nil {
		return u.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepo) FindUserByNotifyTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*entity.User, error) {
	args := m.Called(ctx, startTime, endTime)
	if u := args.Get(0); u != nil {
		return u.([]*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// UserUseCase の生成ヘルパー
func setupUserUseCaseTest() (*MockUserRepo, usecase.UserUseCase) {
	mockRepo := new(MockUserRepo)
	usecase := usecase.NewUserUseCase(mockRepo)
	return mockRepo, usecase
}

// Create のテスト
func TestUserUseCase_Create_Success(t *testing.T) {
	mockRepo, uuc := setupUserUseCaseTest()
	ctx := context.Background()
	now := time.Now().In(utils.JST)

	user := &entity.User{
		LINEUserID:     "U123",
		SelectedAreaID: "1",
		NotifyTime:     time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, utils.JST),
	}

	// モック設定
	mockRepo.On("CreateUser", ctx, user).Return(user, nil)

	created, err := uuc.Create(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, user, created)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Create_InvalidLINEUserID(t *testing.T) {
	_, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	user := &entity.User{
		LINEUserID:     "",
		SelectedAreaID: "1",
		NotifyTime:     time.Now().In(utils.JST),
	}

	created, err := uuc.Create(ctx, user)
	assert.Nil(t, created)
	assert.EqualError(t, err, "LINEUserID is required")
}

// GetByID のテスト
func TestUserUseCase_GetByID_Success(t *testing.T) {
	mockRepo, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	expectedUser := &entity.User{
		ID:             1,
		LINEUserID:     "U123",
		SelectedAreaID: "1",
		NotifyTime:     time.Now().In(utils.JST),
		IsActive:       true,
	}

	mockRepo.On("FindUserByID", ctx, 1).Return(expectedUser, nil)

	user, err := uuc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetByID_InvalidID(t *testing.T) {
	_, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	user, err := uuc.GetByID(ctx, 0)
	assert.Nil(t, user)
	assert.EqualError(t, err, "invalid user id")
}

// GetByLINEID のテスト
func TestUserUseCase_GetByLINEID_Success(t *testing.T) {
	mockRepo, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	expectedUser := &entity.User{
		ID:             2,
		LINEUserID:     "U456",
		SelectedAreaID: "2",
		NotifyTime:     time.Now().In(utils.JST),
		IsActive:       true,
	}

	mockRepo.On("FindUserByLINEUserID", ctx, "U456").Return(expectedUser, nil)

	user, err := uuc.GetByLINEID(ctx, "U456")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetByLINEID_EmptyID(t *testing.T) {
	_, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	user, err := uuc.GetByLINEID(ctx, "")
	assert.Nil(t, user)
	assert.EqualError(t, err, "LINEUserID is required")
}

// Update のテスト
func TestUserUseCase_Update_Success(t *testing.T) {
	mockRepo, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	user := &entity.User{
		ID:             1,
		SelectedAreaID: "3",
		NotifyTime:     time.Now().In(utils.JST),
		IsActive:       false,
	}

	mockRepo.On("UpdateUser", ctx, user).Return(nil)

	err := uuc.Update(ctx, user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Update_InvalidID(t *testing.T) {
	_, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	user := &entity.User{
		ID:             0,
		SelectedAreaID: "3",
		NotifyTime:     time.Now().In(utils.JST),
		IsActive:       false,
	}

	err := uuc.Update(ctx, user)
	assert.EqualError(t, err, "invalid user id")
}

// Delete のテスト
func TestUserUseCase_Delete_Success(t *testing.T) {
	mockRepo, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	mockRepo.On("DeleteUser", ctx, 1).Return(nil)

	err := uuc.Delete(ctx, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Delete_InvalidID(t *testing.T) {
	_, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	err := uuc.Delete(ctx, 0)
	assert.EqualError(t, err, "invalid user id")
}

func TestUserUseCase_Delete_Error(t *testing.T) {
	mockRepo, uuc := setupUserUseCaseTest()
	ctx := context.Background()

	mockRepo.On("DeleteUser", ctx, 2).Return(errors.New("delete failed"))

	err := uuc.Delete(ctx, 2)
	assert.EqualError(t, err, "delete failed")
	mockRepo.AssertExpectations(t)
}
