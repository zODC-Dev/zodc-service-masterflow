package interfaces

import (
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
)

type IFormController interface {
	Create(ctx echo.Context) error
	FindAll() (*[]entities.Form, error)
}

type IFormService interface {
	Create(req *requests.FormCreateRequest) error
	FindAll() (*[]entities.Form, error)
}

type IFormRepository interface {
	Create(req *requests.FormCreateRequest) error
	Delete(form *entities.Form) error
	FindAll() (*[]entities.Form, error)
	Find(form *entities.Form) error
}
