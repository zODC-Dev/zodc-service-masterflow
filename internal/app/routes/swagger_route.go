package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/zODC-Dev/zodc-service-masterflow/docs"
)

func RegisterSwaggerRoute(e *echo.Echo) {
	e.GET("/api/v1/docs/*", echoSwagger.WrapHandler)
}
