package interfaces

import (
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
)

type IFormController interface {
	Create(ctx echo.Context) error
}

type IFormService interface {
	Create(form *requests.FormCreateRequest) error
}

type IFormRepository interface {
	Create(form *requests.FormCreateRequest) error
}
