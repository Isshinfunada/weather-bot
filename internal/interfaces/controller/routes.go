package controller

import (
	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, userUC usecase.UserUsecase, areaUC usecase.AreaUseCase) {
	userCtrl := NewUserController(userUC)
	areaCtrl := NewAreaController(areaUC)

	// User
	e.POST("/api/users", userCtrl.Create)                          //Create
	e.GET("/api/users/:id", userCtrl.GetByID)                      // Read(ByID)
	e.GET("/api/users/line/:lineUserid", userCtrl.GetByLINEUserID) // Read(ByLINEID)
	e.PUT("/api/users/:id", userCtrl.Update)                       // Update
	e.DELETE("/api/users/:id", userCtrl.Delete)                    //Delete

	// Area
	e.GET("/api/areas/:class20_id", areaCtrl.GetHierarchy) //Read
}
