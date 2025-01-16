package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/controllers"
)

func UtilRoute(group *echo.Group) {
	utilGroup := group.Group("/utils")
	{
		utilGroup.POST("/excel/extract", controllers.ExcelExtract)
	}
}
