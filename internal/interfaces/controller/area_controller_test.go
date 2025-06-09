package controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/controller"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAreaUseCase は usecase.AreaUseCase インターフェースのモック実装
type MockAreaUseCase struct {
	mock.Mock
}

func (m *MockAreaUseCase) SearchCityCandidates(ctx context.Context, query string) ([]*entity.HierarchyArea, error) {
	args := m.Called(ctx, query)
	if c := args.Get(0); c != nil {
		return c.([]*entity.HierarchyArea), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAreaUseCase) GetHierarchy(ctx context.Context, class20ID string) (*entity.HierarchyArea, error) {
	args := m.Called(ctx, class20ID)
	if h := args.Get(0); h != nil {
		return h.(*entity.HierarchyArea), args.Error(1)
	}
	return nil, args.Error(1)
}

// AreaController 用のモックユースケースとコントローラーのセットアップ
func setupAreaControllerTest() (*MockAreaUseCase, *controller.AreaController, echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	// クエリパラメータに対応するため、パス内のIDは仮に設定
	req := httptest.NewRequest(http.MethodGet, "/api/areas/1234567", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(MockAreaUseCase)
	controller := controller.NewAreaController(mockUC)
	return mockUC, controller, c, rec
}

// 無効な class20_id パラメータによるバリデーションエラーのテスト
// ※クラスIDが文字列の場合、型変換エラーは発生しないため、このケースは不要。
// 代わりに、ユースケース側で不正な文字列に対する処理を検証できます。

// Areaが見つからない場合のテスト
func TestAreaController_GetHierarchy_NotFound(t *testing.T) {
	mockUC, ctrl, c, rec := setupAreaControllerTest()

	// URLパラメータ設定
	c.SetParamNames("class20_id")
	c.SetParamValues("1234567")

	ctx := c.Request().Context()
	validIDStr := "1234567"
	notFoundErr := errors.New("not found for class20 id=1234567")

	mockUC.
		On("GetHierarchy", ctx, validIDStr).
		Return(nil, notFoundErr)

	err := ctrl.GetHierarchy(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, notFoundErr.Error(), resp["error"])

	mockUC.AssertExpectations(t)
}

// 正常系のテスト
func TestAreaController_GetHierarchy_Success(t *testing.T) {
	mockUC, ctrl, c, rec := setupAreaControllerTest()

	// URLパラメータ設定
	c.SetParamNames("class20_id")
	c.SetParamValues("1234567")

	ctx := c.Request().Context()
	validIDStr := "1234567"
	expectedHierarchy := &entity.HierarchyArea{
		Class20: &entity.AreaClass20{ID: "1234567", Name: "TestArea"},
		// 必要に応じて他の階層のデータを設定
	}

	mockUC.
		On("GetHierarchy", ctx, validIDStr).
		Return(expectedHierarchy, nil)

	err := ctrl.GetHierarchy(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp entity.HierarchyArea
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedHierarchy.Class20.ID, resp.Class20.ID)
	assert.Equal(t, expectedHierarchy.Class20.Name, resp.Class20.Name)

	mockUC.AssertExpectations(t)
}
