package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/interfaces"
)

type formControllerImpl struct {
	formService interfaces.IFormService
}

func NewFormController(formService interfaces.IFormService) *formControllerImpl {
	return &formControllerImpl{
		formService: formService,
	}
}

func (c *formControllerImpl) Create(ctx echo.Context) error {
	req := new(requests.FormCreateRequest)

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

func (c *formControllerImpl) FindAll(ctx echo.Context) error {
	forms, err := c.formService.FindAll()
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, err.Error())
	}

	return ctx.JSON(http.StatusCreated, forms)
}
