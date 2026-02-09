package router

import (
	"net/http"
	"scadulDataMono/infra/echo_http/middleware"
	"scadulDataMono/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterTermRoutes(e *echo.Echo, uc *usecase.Term) {
	g := e.Group("/term")

	g.POST("", func(c echo.Context) error {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		id, err := uc.Create(req.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"id": id})
	}, middleware.Permit(0))

	g.GET("", func(c echo.Context) error {
		list, err := uc.Listing()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, list)
	}, middleware.Permit(0, 1))

	g.PUT("/:id", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		var req struct {
			Name string `json:"name"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		updated, err := uc.Update(uint(id), req.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, updated)
	}, middleware.Permit(0))

	g.DELETE("/:id", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		if err := uc.Delete(uint(id)); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))
}
