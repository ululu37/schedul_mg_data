package router

import (
	"net/http"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/echo_http/middleware"
	"scadulDataMono/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterStudentRoutes(e *echo.Echo, uc *usecase.StudentMg) {
	g := e.Group("/student")

	g.POST("", func(c echo.Context) error {
		var req struct {
			Name         string
			CurriculumID uint
			Year         int
			ClassroomID  uint
			Username     string
			Password     string
			Role         int
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		id, err := uc.Create(req.Name, req.Year, req.Username, req.Password, req.Role)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"id": id})
	}, middleware.Permit(0))

	g.GET("", func(c echo.Context) error {
		search := c.QueryParam("search")
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("perpage"))
		list, count, err := uc.Listing(search, page, perPage)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"data": list, "count": count})
	}, middleware.Permit(0, 1))

	g.GET("/:id", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		student, err := uc.GetByID(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, student)
	}, middleware.Permit(0, 1))

	g.PUT("/:id", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		var req entities.Student
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		updated, err := uc.Update(uint(id), &req)
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
