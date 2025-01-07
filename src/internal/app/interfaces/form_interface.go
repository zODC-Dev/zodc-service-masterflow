package interfaces

import (
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/models"
)

type IFormController interface {
	Create(ctx echo.Context) error
}

type IFormService interface {
	Create(form *models.FormCreateRequest) error
}

type IFormRepository interface {
	Create(form *models.FormCreateRequest) error
}
