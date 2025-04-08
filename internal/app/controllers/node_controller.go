package controllers

import (
	"net/http"
	"strconv"

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

func (c *NodeController) ReassignNode(e echo.Context) error {
	ctx := e.Request().Context()

	nodeId := e.Param("id")
	userId := e.Param("userId")

	if nodeId == "" || userId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and userId are required"})
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid userId"})
	}

	if err := c.nodeService.ReassignNode(ctx, nodeId, int32(userIdInt)); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node re-assigned successfully"})
}

// func (c *NodeController) SubmitNodeForm(e echo.Context) error {
// 	ctx := e.Request().Context()

// 	nodeId := e.Param("id")
// 	formId := e.Param("formId")

// 	req := new(requests.SubmitNodeFormRequest)
// 	if err := e.Bind(req); err != nil {
// 		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
// 	}

// 	if nodeId == "" || formId == "" {
// 		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and formId are required"})
// 	}

// 	if err := c.nodeService.SubmitNodeForm(ctx, nodeId, formId, req); err != nil {
// 		return e.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	return e.JSON(http.StatusOK, map[string]string{"message": "Node form submitted successfully"})
// }

// func (c *NodeController) ApproveNodeForm(e echo.Context) error {
// 	ctx := e.Request().Context()

// 	nodeId := e.Param("id")
// 	formId := e.Param("formId")

// 	if nodeId == "" || formId == "" {
// 		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and formId are required"})
// 	}

// 	if err := c.nodeService.ApproveNodeForm(ctx, nodeId, formId); err != nil {
// 		return e.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	return e.JSON(http.StatusOK, map[string]string{"message": "Node form approved successfully"})
// }

// func (c *NodeController) RejectNodeForm(e echo.Context) error {
// 	ctx := e.Request().Context()

// 	nodeId := e.Param("id")
// 	formId := e.Param("formId")

// 	if nodeId == "" || formId == "" {
// 		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and formId are required"})
// 	}

// 	if err := c.nodeService.RejectNodeForm(ctx, nodeId, formId); err != nil {
// 		return e.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	return e.JSON(http.StatusOK, map[string]string{"message": "Node form rejected successfully"})
// }
