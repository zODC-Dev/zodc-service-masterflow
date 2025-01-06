package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/models"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/services"
)

type IFormController interface {
	Create(ctx echo.Context) error
}

type formControllerImpl struct {
	formService services.IFormService
}

func NewFormController(formService services.IFormService) *formControllerImpl {
	return &formControllerImpl{
		formService: formService,
	}
}

func (c *formControllerImpl) Create(ctx echo.Context) error {
	req := new(models.FormCreateRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.formService.Create(req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{
		"message": "Form created successfully",
	})

}
