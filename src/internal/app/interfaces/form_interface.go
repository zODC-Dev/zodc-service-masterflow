package interfaces

import (
	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/types"
)

type FormController interface {
	FindAll(ctx echo.Context) error
	Create(ctx echo.Context) error
}

type FormService interface {
	FindAll() (*[]responses.FormFindAll, error)
	Create(*requests.FormCreate) error
}

type FormRepository interface {
	FindAll() (*[]types.FormWithFields, error)
	CreateForm(form model.Forms, formFields []model.FormFields) error
}
