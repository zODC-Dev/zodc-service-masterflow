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

// CreateFormTemplate godoc
// @Summary      Create a new form template
// @Description  Adds a new form template to the system based on the provided data.
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        formTemplate body requests.FormTemplateCreate true "Form Template Creation Request"
// @Success      201 {object} map[string]string "message: Form created successfully"
// @Failure      400 {object} string "Error message for bad request"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /forms/templates [post]
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

// FindAllFormTemplate godoc
// @Summary      Find all form templates
// @Description  Retrieves a list of form templates, optionally filtered by category ID and search query.
// @Tags         Forms
// @Produce      json
// @Param        categoryId query string false "Filter by category ID"
// @Param        search query string false "Search query for form template names or descriptions"
// @Success      200 {array} responses.FormTemplateFindAll // Assuming responses.FormTemplateResponse exists
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /forms/templates [get]
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

// FindAllFormTemplateFieldsByFormTemplateId godoc
// @Summary      Find all fields for a specific form template
// @Description  Retrieves all fields associated with a given form template ID.
// @Tags         Forms
// @Produce      json
// @Param        formTemplateId path int true "Form Template ID"
// @Success      200 {array} responses.FormTemplateFieldsFindAll // Assuming responses.FormTemplateFieldResponse exists
// @Failure      400 {object} string "Error message for invalid form template ID"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /forms/templates/{formTemplateId}/fields [get]
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

func (c *FormController) FindOneFormTemplateDetailByFormTemplateId(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplateId := e.Param("formTemplateId")
	if formTemplateId == "" {
		return e.JSON(http.StatusBadRequest, "Form template ID is required")
	}

	formTemplateIdInt, err := strconv.Atoi(formTemplateId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Form template ID is not a valid integer")
	}

	formTemplateDetails, err := c.formService.FindOneFormTemplateDetailByFormTemplateId(ctx, int32(formTemplateIdInt))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, formTemplateDetails)

}

func (c *FormController) UpdateFormTemplate(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplateId := e.Param("formTemplateId")
	if formTemplateId == "" {
		return e.JSON(http.StatusBadRequest, "Form template ID is required")
	}

	formTemplateIdInt, err := strconv.Atoi(formTemplateId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Form template ID is not a valid integer")
	}

	req := new(requests.FormTemplateUpdate)
	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.formService.UpdateFormTemplate(ctx, req, int32(formTemplateIdInt)); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "Form template updated successfully",
	})
}
