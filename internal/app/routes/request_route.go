package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

func RequestRoute(group *echo.Group, db *sql.DB) {

	requestRepo := repositories.NewRequestRepository()

	userApi := externals.NewUserAPI()

	requestService := services.NewRequestService(db, requestRepo, userApi)

	requestController := controllers.NewRequestController(requestService)

	requestRoute := group.Group("/requests")
	{
		requestRoute.GET("", requestController.FindAllRequest)
		requestRoute.GET("/overview", requestController.GetRequestOverview)
	}

}
