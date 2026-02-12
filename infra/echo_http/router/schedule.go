package router

import (
	"net/http"
	"scadulDataMono/usecase"

	"github.com/labstack/echo/v4"
)

func RegisterScheduleRoutes(e *echo.Echo, scheduleUsecase *usecase.ScheduleUsecase) {
	e.POST("/schedule/generate", func(c echo.Context) error {
		res, err := scheduleUsecase.GenerateSchedule()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, res)
	})
}
