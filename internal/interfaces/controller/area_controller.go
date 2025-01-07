package controller

import (
	"net/http"

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
	class20ID := c.Param("class20_id")

	ctx := c.Request().Context()
	hierarchy, err := ctrl.areaUC.GetHierarchy(ctx, class20ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, hierarchy)
}
