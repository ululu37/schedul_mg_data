package router

import (
	"net/http"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/echo_http/middleware"
	"scadulDataMono/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterTeacherRoutes(e *echo.Echo, uc *usecase.TeacherMg, tEverlute *usecase.TeacherEverlute) {
	g := e.Group("/teacher")

	g.GET("/:id", func(c echo.Context) error {
		teacherID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid teacher id")
		}
		teacher, err := uc.GetByID(uint(teacherID))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, teacher)
	}, middleware.Permit(0, 1))

	g.POST("", func(c echo.Context) error {
		var req struct {
			Name     string `json:"name"`
			Resume   string `json:"resume"`
			Username string `json:"username"`
			Password string `json:"password"`
			Role     int    `json:"role"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		id, err := uc.Create(req.Name, req.Resume, req.Username, req.Password, req.Role)
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

	g.PUT("/:id", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		var req struct {
			Name   string `json:"name"`
			Resume string `json:"resume"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		updated, err := uc.Update(uint(id), req.Name, req.Resume)
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

	g.GET("/:id/mysubject", func(c echo.Context) error {
		teacherID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("perpage"))
		list, count, err := uc.GetMySubject(uint(teacherID), 1, page, perPage)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"data": list, "count": count})
	}, middleware.Permit(0, 1))

	g.POST("/:id/mysubject", func(c echo.Context) error {
		teacherID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var req []struct {
			SubjectID  uint `json:"subject_id"`
			Preference int  `json:"preference"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var mySubjects []entities.TeacherMySubject
		for _, req := range req {
			newMysubject := entities.TeacherMySubject{
				SubjectID:  req.SubjectID,
				Preference: req.Preference,
			}
			mySubjects = append(mySubjects, newMysubject)
		}
		err := uc.AddMySubject(uint(teacherID), mySubjects)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.DELETE("/:id/mysubject", func(c echo.Context) error {
		var req struct {
			IDs []uint `json:"ids"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := uc.RemoveMySubject(req.IDs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.POST("/aieverlute", func(c echo.Context) error {
		err := tEverlute.Everlute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))
}
