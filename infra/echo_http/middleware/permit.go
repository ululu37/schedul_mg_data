package middleware

import (
	"net/http"
	jwthast "scadulDataMono/infra/jwt_hast"

	"github.com/labstack/echo/v4"
)

func Permit(roles ...int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from cookie or header
			tokenString := ""
			cookie, err := c.Cookie("token")
			//fmt.Println("cookie:", cookie)
			if err == nil {
				tokenString = cookie.Value
			} else {
				// Check Authorization Header
				authHeader := c.Request().Header.Get("Authorization")
				if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
					tokenString = authHeader[7:]
				}
			}

			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, "Missing token")
			}

			// Parse token
			claims, err := jwthast.ParseToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "Invalid token")
			}

			// Check role
			isAllowed := false
			for _, role := range roles {
				if claims.Payload.Role == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				return c.JSON(http.StatusForbidden, "Access denied")
			}

			// Set user payload in context
			c.Set("user", claims.Payload)

			return next(c)
		}
	}
}
