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
		return e.JSON(http.StatusBadRequest, err.Error())
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
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusCreated, formTemplates)
}

// FindOneFormTemplateDetailByFormTemplateId godoc
// @Summary      Get form template detail by ID
// @Description  Retrieves detailed information of a specific form template using its ID.
// @Tags         Forms
// @Produce      json
// @Param        formTemplateId path int true "Form Template ID"
// @Success      200 {object} responses.FormTemplateDetails
// @Failure      400 {object} map[string]string "error: Invalid formTemplateId or service error"
// @Router       /forms/templates/{formTemplateId}/details [get]
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
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, formTemplateDetails)

}

// UpdateFormTemplate godoc
// @Summary      Update a form template
// @Description  Updates the details of a specific form template by its ID.
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        formTemplateId path int true "Form Template ID"
// @Param        formTemplate body requests.FormTemplateUpdate true "Form Template update payload"
// @Success      200 {object} map[string]string "message: Form template updated successfully"
// @Failure      400 {object} map[string]string "error: Invalid ID, bad request, or service error"
// @Router       /forms/templates/{formTemplateId} [put]
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
		return e.JSON(http.StatusBadRequest, err.Error())
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
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "Form template config updated successfully",
	})
}

// FindOneFormDataByFormDataId godoc
// @Summary      Get form data by ID
// @Description  Retrieves detailed data of a specific submitted form by its ID.
// @Tags         FormData
// @Produce      json
// @Param        formDataId path int true "Form Data ID"
// @Success      200 {object} responses.JiraFormDetailResponse
// @Failure      400 {object} map[string]string "error: Invalid ID or service error"
// @Router       /forms/data/{formDataId} [get]
func (c *FormController) FindOneFormDataByFormDataId(e echo.Context) error {
	ctx := e.Request().Context()

	formDataId := e.Param("dataId")
	if formDataId == "" {
		return e.JSON(http.StatusBadRequest, "Form data ID is required")
	}

	formData, err := c.formService.FindOneFormDataByFormDataId(ctx, formDataId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, formData)
}

// GetEditProfileFormTemplate godoc
// @Summary      Get edit profile form template
// @Description  Retrieves the form template used for editing a user's profile.
// @Tags         Forms
// @Produce      json
// @Success      200 {object} responses.FormTemplateDetails
// @Failure      400 {object} map[string]string "error: Service error or unable to retrieve the form template"
// @Router       /forms/templates/edit-profile [get]
func (c *FormController) GetEditProfileFormTemplate(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplate, err := c.formService.GetEditProfileFormTemplate(ctx)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, formTemplate)
}

// GetPerformanceEvaluateFormTemplate godoc
// @Summary      Get performance evaluation form template
// @Description  Retrieves the form template used for performance evaluations.
// @Tags         Forms
// @Produce      json
// @Success      200 {object} responses.FormTemplateDetails
// @Failure      400 {object} map[string]string "error: Service error or unable to retrieve the form template"
// @Router       /forms/templates/performance-evaluate [get]
func (c *FormController) GetPerformanceEvaluateFormTemplate(e echo.Context) error {
	ctx := e.Request().Context()

	formTemplate, err := c.formService.GetPerformanceEvaluateFormTemplate(ctx)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, formTemplate)
}
