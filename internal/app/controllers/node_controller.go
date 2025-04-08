package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/middlewares"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

type NodeController struct {
	nodeService *services.NodeService
}

func NewNodeController(nodeService *services.NodeService) *NodeController {
	return &NodeController{
		nodeService: nodeService,
	}
}

func (c *NodeController) CompleteNode(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	if nodeId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId is required"})
	}

	if err := c.nodeService.CompleteNodeHandler(ctx, nodeId, userId); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node completed successfully"})
}

func (c *NodeController) StartNode(e echo.Context) error {
	ctx := e.Request().Context()

	nodeId := e.Param("id")
	if err := c.nodeService.StartNodeHandler(ctx, nodeId); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node started successfully"})
}

func (c *NodeController) GetNodeFormWithPermission(e echo.Context) error {
	ctx := e.Request().Context()

	nodeId := e.Param("id")
	permission := e.Param("permission")

	if nodeId == "" || permission == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and permission are required"})
	}

	form, err := c.nodeService.GetNodeFormWithPermission(ctx, nodeId, permission)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, form)
}

func (c *NodeController) GetNodeJiraForm(e echo.Context) error {
	ctx := e.Request().Context()

	nodeId := e.Param("id")
	if nodeId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId is required"})
	}

	jiraForm, err := c.nodeService.GetNodeJiraForm(ctx, nodeId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, jiraForm)
}

func (c *NodeController) GetNodeTaskDetail(e echo.Context) error {
	ctx := e.Request().Context()

	nodeId := e.Param("id")
	if nodeId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId is required"})
	}

	tasks, err := c.nodeService.GetNodeTaskDetail(ctx, nodeId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, tasks)
}
