package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

type FormController struct {
	formService *services.FormService
}

func NewFormController(formService *services.FormService) *FormController {
	return &FormController{
		formService: formService,
	}
}

func (c *FormController) CreateFormTemplate(e echo.Context) error {
	ctx := e.Request().Context()

	req := new(requests.FormTemplateCreate)

	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.formService.CreateFormTemplate(ctx, req); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, map[string]string{
		"message": "Form created successfully",
	})
}

func (c *FormController) FindAllFormTemplate(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplates, err := c.formService.FindAllFormTemplate(ctx)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, formTemplates)
}
