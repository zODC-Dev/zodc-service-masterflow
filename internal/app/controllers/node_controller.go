package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
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

// CompleteNode godoc
// @Summary      Complete a node
// @Description  Marks a specific node as completed by the logged-in user.
// @Tags         Nodes
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path string true "Node ID"
// @Success      200 {object} map[string]string "message: Node completed successfully"
// @Failure      400 {object} map[string]string "error: Error message for bad request (e.g., missing ID, invalid ID, service error)"
// @Router       /nodes/{id}/complete [patch] // Assuming PATCH, could be POST
func (c *NodeController) CompleteNode(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	if nodeId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId is required"})
	}

	if err := c.nodeService.CompleteNodeHandler(ctx, nodeId, int32(userId)); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node completed successfully"})
}

// StartNode godoc
// @Summary      Start a node
// @Description  Marks a specific node as started.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Success      200 {object} map[string]string "message: Node started successfully"
// @Failure      400 {object} string "Error message for bad request (e.g., missing ID, service error)"
// @Router       /nodes/{id}/start [patch] // Assuming PATCH, could be POST
func (c *NodeController) StartNode(e echo.Context) error {
	ctx := e.Request().Context()

	nodeId := e.Param("id")
	if err := c.nodeService.StartNodeHandler(ctx, nodeId); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node started successfully"})
}

// GetNodeFormWithPermission godoc
// @Summary      Get node form based on permission
// @Description  Retrieves the form associated with a node, considering the user's permission level.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Param        permission path string true "Permission level (e.g., 'read', 'write')"
// @Success      200 {object} responses.NodeFormResponse // Assuming responses.NodeFormResponse exists
// @Failure      400 {object} map[string]string "error: Error message for bad request (e.g., missing params, service error)"
// @Router       /nodes/{id}/form/{permission} [get]
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

// GetNodeJiraForm godoc
// @Summary      Get Jira form for a node
// @Description  Retrieves the Jira-specific form details associated with a node.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Success      200 {object} responses.JiraFormDetailResponse // Assuming responses.NodeJiraFormResponse exists
// @Failure      400 {object} map[string]string "error: Error message for bad request (e.g., missing ID, service error)"
// @Router       /nodes/{id}/jira-form [get]
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

func (c *NodeController) SubmitNodeForm(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	formId := e.Param("formId")

	req := new([]requests.SubmitNodeFormRequest)
	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if nodeId == "" || formId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and formId are required"})
	}

	if err := c.nodeService.SubmitNodeForm(ctx, userId, nodeId, formId, req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node form submitted successfully"})
}

func (c *NodeController) ApproveNodeForm(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	formId := e.Param("formId")

	if nodeId == "" || formId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and formId are required"})
	}

	if err := c.nodeService.ApproveNodeForm(ctx, nodeId, formId, int32(userId)); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node form approved successfully"})
}

func (c *NodeController) RejectNodeForm(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	formId := e.Param("formId")

	if nodeId == "" || formId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and formId are required"})
	}

	if err := c.nodeService.RejectNodeForm(ctx, nodeId, formId, int32(userId)); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node form rejected successfully"})
}

func (c *NodeController) ApproveNode(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	if nodeId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId is required"})
	}

	if err := c.nodeService.ApproveNode(ctx, int32(userId), nodeId); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node approved successfully"})
}

func (c *NodeController) RejectNode(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	if nodeId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId is required"})
	}

	if err := c.nodeService.RejectNode(ctx, int32(userId), nodeId); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node rejected successfully"})
}

// GetNodeTaskDetail godoc
// @Summary      Get task details for a node
// @Description  Retrieves detailed information about the tasks associated with a specific node.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Success      200 {array} responses.TaskDetail // Assuming responses.TaskResponse exists
// @Failure      400 {object} map[string]string "error: Error message for bad request (e.g., missing ID, service error)"
// @Router       /nodes/{id}/tasks [get]
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
