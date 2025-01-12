package interfaces

import (
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
)

type IFormController interface {
	Create(ctx echo.Context) error
	FindAll(ctx echo.Context) error
	Delete(form *entities.Form) error
}

type IFormService interface {
	Create(req *requests.FormCreateRequest) error
	FindAll() (*[]responses.FormResponse, error)
	Delete(form *entities.Form) error
	FindById(id string) (*entities.Form, error)
}

type IFormRepository interface {
	Create(req *entities.Form) error
	FindAll() (*[]entities.Form, error)
	Delete(form *entities.Form) error
	FindById(id string) (*entities.Form, error)
}
