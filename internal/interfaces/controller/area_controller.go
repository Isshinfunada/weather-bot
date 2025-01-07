package controller

import (
	"net/http"
	"strconv"

	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AreaController struct {
	areaUC usecase.AreaUseCase
}

func NewAreaController(aUC usecase.AreaUseCase) *AreaController {
	return &AreaController{areaUC: aUC}
}

func (ctrl *AreaController) GetHierarchy(c echo.Context) error {
	class20ID, err := strconv.Atoi(c.Param("class20_id"))
	if err != nil || class20ID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "class20_id is requiredand must be a valid integer"})
	}

	ctx := c.Request().Context()
	hierarchy, err := ctrl.areaUC.GetHierarchy(ctx, class20ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, hierarchy)
}
