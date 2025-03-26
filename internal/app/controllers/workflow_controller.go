package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
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

func (c *WorkflowController) StartWorkflow(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	req := new(requests.StartWorkflow)
	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.workflowService.StartWorkflowHandler(ctx, *req, userId); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, nil)
}
