package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/middlewares"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

type WorkflowController struct {
	workflowService *services.WorkflowService
}

func NewWorkflowController(workflowService *services.WorkflowService) *WorkflowController {
	return &WorkflowController{
		workflowService: workflowService,
	}
}

// CreateWorkflow godoc
// @Summary      Create a new workflow
// @Description  Adds a new workflow definition based on the provided data. Requires user authentication.
// @Tags         Workflows
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        workflow body requests.CreateWorkflow true "Workflow Creation Request"
// @Success      201 {object} map[string]string "message: Workflow created successfully"
// @Failure      400 {object} string "Error message for bad request body"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /workflows [post]
func (c *WorkflowController) CreateWorkflow(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	req := new(requests.CreateWorkflow)

	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.workflowService.CreateWorkflowHandler(ctx, req, userId); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, map[string]string{
		"message": "Workflow created successfully",
	})
}

// FindAllWorkflow godoc
// @Summary      Find all workflows
// @Description  Retrieves a list of workflows based on various filter criteria for the logged-in user.
// @Tags         Workflows
// @Produce      json
// @Security     ApiKeyAuth
// @Param        categoryId query string false "Filter by category ID"
// @Param        search query string false "Search term for workflows"
// @Param        type query string false "Filter by workflow type"
// @Param        projectKey query string false "Filter by project key"
// @Param        hasSubWorkflow query bool false "Filter by whether it has sub-workflows"
// @Param        isArchived query bool false "Filter by archival status"
// @Success      200 {array} responses.WorkflowResponse // Assuming responses.WorkflowResponse exists
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /workflows [get]
func (c *WorkflowController) FindAllWorkflow(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	workflowTemplateQueryParams := queryparams.WorkflowQueryParam{
		CategoryID:     e.QueryParam("categoryId"),
		Search:         e.QueryParam("search"),
		Type:           e.QueryParam("type"),
		ProjectKey:     e.QueryParam("projectKey"),
		HasSubWorkflow: e.QueryParam("hasSubWorkflow"),
		IsArchived:     e.QueryParam("isArchived"),
	}

	workflows, err := c.workflowService.FindAllWorkflowHandler(ctx, workflowTemplateQueryParams, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, workflows)
}

// FindOneWorkflowDetail godoc
// @Summary      Find workflow detail by version ID
// @Description  Retrieves detailed information for a specific workflow version ID.
// @Tags         Workflows
// @Produce      json
// @Param        id path int true "Workflow Version ID"
// @Success      200 {object} responses.WorkflowDetailResponse // Assuming responses.WorkflowDetailResponse exists
// @Failure      500 {object} string "Error message for invalid ID or internal server error" // Combined 400/500 based on current code
// @Router       /workflows/{id} [get] // Assuming ID refers to version ID based on param name
func (c *WorkflowController) FindOneWorkflowDetail(e echo.Context) error {
	ctx := e.Request().Context()

	workflowVersionId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	workflowDetail, err := c.workflowService.FindOneWorkflowDetailHandler(ctx, int32(workflowVersionId))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, workflowDetail)
}

// StartWorkflow godoc
// @Summary      Start a workflow instance
// @Description  Initiates a new instance of a workflow based on the provided request details. Requires user authentication.
// @Tags         Workflows
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        startRequest body requests.StartWorkflow true "Start Workflow Request"
// @Success      200 {object} responses.Response{data=int32} // Assuming appropriate response structure
// @Failure      400 {object} string "Error message for bad request body"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /workflows/start [post]
func (c *WorkflowController) StartWorkflow(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	req := new(requests.StartWorkflow)
	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	requestId, err := c.workflowService.StartWorkflowHandler(ctx, *req, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, responses.Response{
		Message: "Workflow started successfully",
		Data: struct {
			ID    int32  `json:"id"`
			Title string `json:"title"`
		}{
			ID:    requestId,
			Title: req.Title,
		},
	})
}

// ArchiveWorkflow godoc
// @Summary      Archive a workflow
// @Description  Marks a workflow definition as archived based on its ID.
// @Tags         Workflows
// @Produce      json
// @Param        id path int true "Workflow ID"
// @Success      200 {object} map[string]string "message: Workflow archived successfully"
// @Failure      500 {object} string "Error message for invalid ID or internal server error" // Combined 400/500 based on current code
// @Router       /workflows/{id}/archive [patch] // Or PUT/DELETE depending on idempotency semantics
func (c *WorkflowController) ArchiveWorkflow(e echo.Context) error {
	ctx := e.Request().Context()

	workflowId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	err = c.workflowService.ArchiveWorkflowHandler(ctx, int32(workflowId))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "Workflow archived successfully",
	})

}
