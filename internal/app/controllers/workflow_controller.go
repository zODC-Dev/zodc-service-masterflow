package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
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

func (c *WorkflowController) CreateWorkflow(e echo.Context) error {
	ctx := e.Request().Context()

	req := new(requests.WorkflowRequest)

	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.workflowService.CreateWorkFlowHandler(ctx, req); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, map[string]string{
		"message": "Workflow created successfully",
	})
}

func (c *WorkflowController) FindAllWorkflow(e echo.Context) error {
	ctx := e.Request().Context()

	workflowTemplateQueryParams := queryparams.WorkflowQueryParam{
		CategoryID: e.QueryParam("categoryId"),
		Search:     e.QueryParam("search"),
		Type:       e.QueryParam("type"),
	}

	workflows, err := c.workflowService.FindAllWorkflowHandler(ctx, workflowTemplateQueryParams)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, workflows)
}

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

	return e.JSON(http.StatusCreated, workflowDetail)
}
