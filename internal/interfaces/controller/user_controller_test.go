package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/controller"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUsecase は usecase.UserUsecase インターフェースのモック実装
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	if u := args.Get(0); u != nil {
		return u.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) GetByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if u := args.Get(0); u != nil {
		return u.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) GetByLINEID(ctx context.Context, LINEUserID string) (*entity.User, error) {
	args := m.Called(ctx, LINEUserID)
	if u := args.Get(0); u != nil {
		return u.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserUsecase) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) Delete(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// テスト用のヘルパー関数：新しい Echo コンテキストと Recorder を生成
func newTestContext(method, path string, body []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// Create エンドポイントのテスト（正常系）
func TestUserController_Create_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	userCtrl := controller.NewUserController(mockUC)

	reqBody := controller.CreateUserRequest{
		LINEUserID:     "U123",
		SelectedAreaID: "1",
		NotifyTime:     "09:00",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	c, rec := newTestContext(http.MethodPost, "/api/users", bodyBytes)

	now := time.Now()
	expectedUser := &entity.User{
		ID:             1,
		LINEUserID:     "U123",
		SelectedAreaID: 1,
		NotifyTime:     time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC),
		IsActive:       false,
	}

	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(expectedUser, nil)

	if assert.NoError(t, userCtrl.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var respUser entity.User
		err := json.Unmarshal(rec.Body.Bytes(), &respUser)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, respUser.ID)
	}
	mockUC.AssertExpectations(t)
}

// GetByID エンドポイントのテスト（正常系）
func TestUserController_GetByID_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	userCtrl := controller.NewUserController(mockUC)

	// URL パラメータとして id を設定するために Echo コンテキストを作成
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	ctx := context.Background()
	expectedUser := &entity.User{
		ID:             1,
		LINEUserID:     "U123",
		SelectedAreaID: 1,
		NotifyTime:     time.Now(),
		IsActive:       true,
	}

	mockUC.On("GetByID", ctx, 1).Return(expectedUser, nil)

	if assert.NoError(t, userCtrl.GetByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var respUser entity.User
		err := json.Unmarshal(rec.Body.Bytes(), &respUser)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, respUser.ID)
	}
	mockUC.AssertExpectations(t)
}

// GetByLINEUserID エンドポイントのテスト（正常系）
func TestUserController_GetByLINEUserID_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	userCtrl := controller.NewUserController(mockUC)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/line/U123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("lineUserid")
	c.SetParamValues("U123")

	ctx := context.Background()
	expectedUser := &entity.User{
		ID:             2,
		LINEUserID:     "U123",
		SelectedAreaID: 2,
		NotifyTime:     time.Now(),
		IsActive:       true,
	}

	mockUC.On("GetByLINEID", ctx, "U123").Return(expectedUser, nil)

	if assert.NoError(t, userCtrl.GetByLINEUserID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var respUser entity.User
		err := json.Unmarshal(rec.Body.Bytes(), &respUser)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, respUser.ID)
	}
	mockUC.AssertExpectations(t)
}

// Update エンドポイントのテスト（正常系）
func TestUserController_Update_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	userCtrl := controller.NewUserController(mockUC)

	reqBody := controller.UpdateUserRequest{
		SelectedAreaID: "2",
		NotifyTime:     "10:00",
		IsActive:       true,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	ctx := context.Background()

	// UpdateUser は error を返さないケースを設定
	mockUC.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	if assert.NoError(t, userCtrl.Update(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "user updated", resp["message"])
	}
	mockUC.AssertExpectations(t)
}

// Delete エンドポイントのテスト（正常系）
func TestUserController_Delete_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	userCtrl := controller.NewUserController(mockUC)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	ctx := context.Background()

	mockUC.On("Delete", ctx, 1).Return(nil)

	if assert.NoError(t, userCtrl.Delete(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "user deleted", resp["message"])
	}
	mockUC.AssertExpectations(t)
}
