package router

import (
	"fmt"
	"net/http"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/echo_http/middleware"
	"scadulDataMono/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterPreCurriculumRoutes(e *echo.Echo, uc *usecase.PreCurriculum) {
	g := e.Group("/precurriculum")

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
		search := c.QueryParam("search")
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("perpage"))
		list, count, err := uc.Listing(search, page, perPage)
		fmt.Println(list)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"data": list, "count": count})
	}, middleware.Permit(0, 1))

	g.GET("/:id", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		data, err := uc.GetByID(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, data)
	}, middleware.Permit(0, 1))

	g.PUT("/:id", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
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
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		if err := uc.Delete(uint(id)); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.POST("/:id/subject", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var req []struct {
			SubjectName string `json:"subject_name"`
			Credit      int    `json:"credit"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		fmt.Println("id:", id)
		var newSubjectInCurriculum []entities.SubjectInPreCurriculum
		for _, r := range req {
			newSubjectInCurriculum = append(newSubjectInCurriculum, entities.SubjectInPreCurriculum{
				PreCurriculumID: uint(id),
				Subject:         entities.Subject{Name: r.SubjectName},
				Credit:          r.Credit,
			})
		}

		err := uc.CreateSubject(uint(id), newSubjectInCurriculum)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.DELETE("/subject/:id", func(c echo.Context) error {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err := uc.RemoveSubject(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))
}
