package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/controllers"
	"gorm.io/gorm"
)

func UtilRoute(group *echo.Group, db *gorm.DB) {
	utilGroup := group.Group("/utils")
	{
		utilGroup.POST("/excel/extract", controllers.ExcelExtract)
	}
}
