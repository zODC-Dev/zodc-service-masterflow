package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
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

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	if err := c.nodeService.StartNodeHandler(ctx, int32(userId), nodeId); err != nil {
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

// ReassignNode godoc
// @Summary      Reassign a node to a new user
// @Description  Reassigns a specific node to a new user based on node ID and user ID. The current user must have permission to reassign the node.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Param        userId path string true "User ID (new assignee)"
// @Success      200 {object} map[string]string "message: Node re-assigned successfully"
// @Failure      400 {object} map[string]string "error: Invalid node ID, user ID, or service error"
// @Router       /nodes/{id}/reassign/{userId} [put]
func (c *NodeController) ReassignNode(e echo.Context) error {
	ctx := e.Request().Context()

	userIdReq, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	userId := e.Param("userId")

	if nodeId == "" || userId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and userId are required"})
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid userId"})
	}

	if err := c.nodeService.ReassignNode(ctx, nodeId, int32(userIdInt), int32(userIdReq)); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node re-assigned successfully"})
}

// SubmitNodeForm godoc
// @Summary      Submit a form associated with a node
// @Description  Submits form data for a specific node based on node ID and form ID. Requires user authentication and valid data.
// @Tags         Nodes
// @Accept       json
// @Produce      json
// @Param        id path string true "Node ID"
// @Param        formId path string true "Form ID"
// @Param        body body []requests.SubmitNodeFormRequest true "Form submission data"
// @Success      200 {object} map[string]string "message: Node form submitted successfully"
// @Failure      400 {object} map[string]string "error: Invalid node ID, form ID, or submission error"
// @Router       /nodes/{id}/forms/{formId}/submit [post]
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

// EditNodeForm godoc
// @Summary      Edit a form associated with a node
// @Description  Edits a specific form data associated with a node, requiring node ID and form ID.
// @Tags         Nodes
// @Accept       json
// @Produce      json
// @Param        id path string true "Node ID"
// @Param        formDataId path string true "Form Data ID"
// @Param        body body []requests.SubmitNodeFormRequest true "Form data submission payload"
// @Success      200 {object} map[string]string "message: Node form edited successfully"
// @Failure      400 {object} map[string]string "error: Invalid node ID, form ID, or service error"
// @Router       /nodes/{id}/forms/{formDataId}/edit [put]
func (c *NodeController) EditNodeForm(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	nodeId := e.Param("id")
	formId := e.Param("formDataId")
	if nodeId == "" || formId == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "nodeId and formId are required"})
	}

	req := new([]requests.SubmitNodeFormRequest)
	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.nodeService.EditNodeForm(ctx, int32(userId), nodeId, formId, req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Node form edited successfully"})
}

// ApproveNodeForm godoc
// @Summary      Approve a form associated with a node
// @Description  Approves a specific form based on node ID and form ID. The user must have permission to approve the form.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Param        formId path string true "Form ID"
// @Success      200 {object} map[string]string "message: Node form approved successfully"
// @Failure      400 {object} map[string]string "error: Invalid node ID, form ID, or service error"
// @Router       /nodes/{id}/forms/{formId}/approve [put]
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

// RejectNodeForm godoc
// @Summary      Reject a form associated with a node
// @Description  Rejects a specific form based on node ID and form ID. The user must have permission to reject the form.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Param        formId path string true "Form ID"
// @Success      200 {object} map[string]string "message: Node form rejected successfully"
// @Failure      400 {object} map[string]string "error: Invalid node ID, form ID, or service error"
// @Router       /nodes/{id}/forms/{formId}/reject [put]
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

// ApproveNode godoc
// @Summary      Approve a node
// @Description  Approves a specific node based on its ID. The user must have permission to approve the node.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Success      200 {object} map[string]string "message: Node approved successfully"
// @Failure      400 {object} map[string]string "error: Invalid node ID, permission issues, or service error"
// @Router       /nodes/{id}/approve [put]
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

// RejectNode godoc
// @Summary      Reject a node
// @Description  Rejects a specific node based on its ID. The user must have permission to reject the node.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Success      200 {object} map[string]string "message: Node rejected successfully"
// @Failure      400 {object} map[string]string "error: Invalid node ID, permission issues, or service error"
// @Router       /nodes/{id}/reject [put]
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

// GetNodeStoryByAssignee godoc
// @Summary      Get stories assigned to a user on a node
// @Description  Retrieves a list of stories that are assigned to a specific user (assignee) on a given node.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Param        assignee query string true "Assignee ID or username"
// @Success      200 {array} responses.WorkflowResponse
// @Failure      400 {object} map[string]string "error: Error message for bad request (e.g., missing ID, service error)"
// @Router       /nodes}/stories [get]
func (c *NodeController) GetNodeStoryByAssignee(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	stories, err := c.nodeService.GetNodeStoryByAssignee(ctx, int32(userId))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, stories)
}

// GetNodeTaskCount godoc
// @Summary      Get task count for a node
// @Description  Returns the total number of tasks associated with a specific node.
// @Tags         Nodes
// @Produce      json
// @Param        id path string true "Node ID"
// @Success      200 {object} responses.NodeTaskCountResponse
// @Failure      400 {object} map[string]string "error: Error message for bad request (e.g., invalid ID, service error)"
// @Router       /nodes/{id}/tasks/count [get]
func (c *NodeController) GetNodeTaskCount(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	count, err := c.nodeService.GetNodeTaskCount(ctx, int32(userId))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, count)
}

// CreateComment godoc
// @Summary      Create a new comment on a node
// @Description  Adds a comment to a specific node, submitted by the authenticated user.
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id path string true "Node ID"
// @Param        comment body requests.CreateComment true "Comment payload"
// @Success      200 {object} map[string]string "message: Comment created successfully"
// @Failure      400 {object} map[string]string "error: Invalid input or service error"
// @Router       /nodes/{id}/comments [post]
func (c *NodeController) CreateComment(e echo.Context) error {
	ctx := e.Request().Context()

	nodeId := e.Param("id")
	userId, _ := middlewares.GetUserID(e)

	req := new(requests.CreateComment)
	if err := e.Bind(req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.nodeService.CreateComment(ctx, req, nodeId, userId); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Comment created successfully"})
}

// GetAllComments godoc
// @Summary      Get all comments for a node
// @Description  Retrieves all comments associated with a specific node.
// @Tags         Comments
// @Produce      json
// @Param        id path string true "Node ID"
// @Success      200 {array} responses.CommentResponse "Success"
// @Failure      400 {object} map[string]string "error: Invalid node ID or service error"
// @Router       /nodes/{id}/comments [get]
func (c *NodeController) GetAllComments(e echo.Context) error {
	ctx := e.Request().Context()

	nodeId := e.Param("id")

	comments, err := c.nodeService.GetAllComments(ctx, nodeId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, responses.Response{
		Message: "Success",
		Data:    comments,
	})
}
