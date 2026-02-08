package router

import (
	"net/http"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/echo_http/middleware"
	"scadulDataMono/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterScadulStudentRoutes(e *echo.Echo, uc *usecase.ScadulStudentMg) {
	g := e.Group("/scadulstudent")

	g.POST("", func(c echo.Context) error {
		var req entities.ScadulStudent
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		id, err := uc.Create(&req)
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
		scadul, err := uc.GetByID(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, scadul)
	}, middleware.Permit(0, 1))

	g.PUT("/:id", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		var req entities.ScadulStudent
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

	g.POST("/:id/subject", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		var req []entities.SubjectInScadulStudent
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		for i := range req {
			req[i].ScadulStudentID = uint(id)
		}
		if err := uc.AddSubjects(req); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.DELETE("/subject", func(c echo.Context) error {
		var req struct {
			IDs []uint `json:"ids"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := uc.RemoveSubjects(req.IDs); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))
}
