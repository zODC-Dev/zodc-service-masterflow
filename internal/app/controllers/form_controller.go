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

func (c *FormController) ConfigFormTemplate(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplateId := e.Param("formTemplateId")
	if formTemplateId == "" {
		return e.JSON(http.StatusBadRequest, "Form template ID is required")
	}

	formTemplateIdInt, err := strconv.Atoi(formTemplateId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Form template ID is not a valid integer")
	}

	req := new([][]requests.FormTemplateFieldsCreate)
	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.formService.ConfigFormTemplate(ctx, int32(formTemplateIdInt), req); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "Form template config updated successfully",
	})
}

func (c *FormController) FindOneFormDataByFormDataId(e echo.Context) error {
	ctx := e.Request().Context()

	formDataId := e.Param("formDataId")
	if formDataId == "" {
		return e.JSON(http.StatusBadRequest, "Form data ID is required")
	}

	formData, err := c.formService.FindOneFormDataByFormDataId(ctx, formDataId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, formData)
}

func (c *FormController) GetEditProfileFormTemplate(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplate, err := c.formService.GetEditProfileFormTemplate(ctx)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, formTemplate)
}

func (c *FormController) GetPerformanceEvaluateFormTemplate(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplate, err := c.formService.GetPerformanceEvaluateFormTemplate(ctx)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, formTemplate)
}
