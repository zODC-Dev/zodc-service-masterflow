package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

func ExtractUserMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDStr := c.Request().Header.Get("x-kong-jwt-claim-sub")
			if userIDStr == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "User ID not found in token",
				})
			}

			userID, err := utils.StringToInt32(userIDStr)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid user ID format",
				})
			}

			c.Set("userID", userID)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (int32, error) {
	userID, ok := c.Get("userID").(int32)

	if !ok {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "User ID not found in context")
	}
	return userID, nil
}
