package router

import (
	"net/http"
	dto "scadulDataMono/domain/DTO"
	"scadulDataMono/domain/entities"
	"scadulDataMono/infra/echo_http/middleware"
	"scadulDataMono/usecase"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Echo, uc *usecase.Auth) {
	g := e.Group("/auth")

	g.POST("/login", func(c echo.Context) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		passport, err := uc.Login(req.Username, req.Password)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		// Set Cookie
		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = passport.Token
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.HttpOnly = true
		cookie.Path = "/" // ⭐ สำคัญมาก
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, passport.Payload)
	})

	g.GET("", func(c echo.Context) error {
		search := c.QueryParam("search")
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("perpage"))
		list, count, err := uc.Listing(search, page, perPage)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"data": list, "count": count})
	}, middleware.Permit(0))

	g.PUT("/:id", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		var req entities.Auth
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
		humanType := c.QueryParam("humantype")
		if err := uc.Delete(uint(id), humanType); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, middleware.Permit(0))

	g.GET("/me", func(c echo.Context) error {
		user := c.Get("user").(dto.PayLoad)
		profile, err := uc.GetProfile(user.ID, user.HumanType)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, profile)
	}, middleware.Permit(0, 1))
}
