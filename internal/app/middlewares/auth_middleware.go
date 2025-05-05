package middlewares

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

func ExtractUserMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			if strings.Contains(path, "/forms/data/") || regexp.MustCompile(`^/requests/\d+/complete$`).MatchString(path) || regexp.MustCompile(`^/sprint/\d+/complete$`).MatchString(path) {
				return next(c)
			}

			var userID int32
			var err error

			authHeader := c.Request().Header.Get("authorization")
			if authHeader != "" {
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				tokenString = strings.Trim(tokenString, "\"")

				if tokenString == "" {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "Invalid token format",
					})
				}

				token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(""), nil
				})
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if userIDStr, ok := claims["sub"].(string); ok && userIDStr != "" {
						userID, err = utils.StringToInt32(userIDStr)
						if err != nil {
							return c.JSON(http.StatusBadRequest, map[string]string{
								"error": "Invalid user ID format",
							})
						}
					} else {
						return c.JSON(http.StatusUnauthorized, map[string]string{
							"error": "User ID not found in token",
						})
					}
				} else {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "Invalid token format",
					})
				}
			} else {
				userIDStr := c.Request().Header.Get("x-kong-jwt-claim-sub")
				if userIDStr == "" {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "User ID not found in token",
					})
				}

				userID, err = utils.StringToInt32(userIDStr)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]string{
						"error": "Invalid user ID format",
					})
				}
			}

			c.Set("userID", userID)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (int32, error) {
	userID, ok := c.Get("userID").(int32)

	if !ok {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "User ID not found in context")
	}
	return userID, nil
}
