package controllers

import (
	"net/http"

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

func (c *WorkflowController) Create(e echo.Context) error {
	ctx := e.Request().Context()

	req := new(requests.WorkflowRequest)

	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.workflowService.Create(ctx, req); err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, map[string]string{
		"message": "Workflow created",
	})
}

func (c *WorkflowController) FindAll(e echo.Context) error {
	ctx := e.Request().Context()

	workflowQueryParam := queryparams.WorkflowQueryParam{
		CategoryID: e.QueryParam("categoryId"),
		Type:       e.QueryParam("type"),
		Search:     e.QueryParam("search"),
	}

	workflows, err := c.workflowService.FindAll(ctx, &workflowQueryParam)
	if err != nil {
		return e.JSON(http.StatusBadGateway, err.Error())
	}

	return e.JSON(http.StatusCreated, workflows)
}
