package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUC usecase.UserUsecase
}

func NewUserController(uuc usecase.UserUsecase) *UserController {
	return &UserController{
		userUC: uuc,
	}
}

// CreateUserRequestはユーザー作成時のJSONリクエストボディ
type CreateUserRequest struct {
	LINEUserID     string `json:"lineUserId"`
	SelectedAreaID string `json:"selectedAreaId"`
	NotifyTime     string `json:"notifyTime"`
}

type UpdateUserRequest struct {
	SelectedAreaID string `json:"selectedAreaId"`
	NotifyTime     string `json:"notifyTime"`
	IsActive       bool   `json:"isActive"`
}

// POST /api/users
func (ctrl *UserController) Create(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	selectedAreaID, err := strconv.Atoi(req.SelectedAreaID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid selected area id"})
	}

	notifyTime, err := time.Parse("15:04", req.NotifyTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid notify time"})
	}

	user := &entity.User{
		LINEUserID:     req.LINEUserID,
		SelectedAreaID: selectedAreaID,
		NotifyTime:     notifyTime,
	}

	ctx := c.Request().Context()
	created, err := ctrl.userUC.Create(ctx, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

// GET /api/users/:id
func (ctrl *UserController) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	ctx := c.Request().Context()
	user, err := ctrl.userUC.GetByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
