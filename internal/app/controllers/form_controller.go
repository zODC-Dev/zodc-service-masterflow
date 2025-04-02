package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
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

	categoryID := e.QueryParam("categoryId")
	search := e.QueryParam("search")

	formTemplates, err := c.formService.FindAllFormTemplate(ctx, queryparams.FormQueryParam{
		CategoryID: categoryID,
		Search:     search,
	})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, formTemplates)
}

func (c *FormController) FindAllFormTemplateFieldsByFormTemplateId(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplateId := e.Param("formTemplateId")
	if formTemplateId == "" {
		return e.JSON(http.StatusBadRequest, "Form template ID is required")
	}

	formTemplateIdInt, err := strconv.Atoi(formTemplateId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Form template ID is not a valid integer")
	}

	formTemplateFields, err := c.formService.FindAllFormTemplateFieldsByFormTemplateId(ctx, int32(formTemplateIdInt))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, formTemplateFields)
}
