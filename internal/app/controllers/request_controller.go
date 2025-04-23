package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/middlewares"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

type RequestController struct {
	requestService *services.RequestService
}

func NewRequestController(requestService *services.RequestService) *RequestController {
	return &RequestController{
		requestService: requestService,
	}
}

// FindAllRequest godoc
// @Summary      Find all requests
// @Description  Retrieves a paginated list of requests based on various filter criteria for the logged-in user.
// @Tags         Requests
// @Produce      json
// @Security     ApiKeyAuth
// @Param        search query string false "Search term for requests"
// @Param        projectKey query string false "Filter by project key"
// @Param        status query string false "Filter by request status"
// @Param        sprintId query string false "Filter by sprint ID"
// @Param        workflowType query string false "Filter by workflow type"
// @Param        page query int false "Page number for pagination" default(1)
// @Param        pageSize query int false "Number of items per page" default(10)
// @Success      200 {object} responses.Response{data=responses.RequestResponse} // Assuming appropriate response structure
// @Failure      400 {object} map[string]string "error: Error message for bad request (e.g., service error)"
// @Router       /requests [get]
func (c *RequestController) FindAllRequest(e echo.Context) error {
	ctx := e.Request().Context()

	search := e.QueryParam("search")
	projectKey := e.QueryParam("projectKey")
	status := e.QueryParam("status")
	sprintId := e.QueryParam("sprintId")
	workflowType := e.QueryParam("workflowType")

	userId, _ := middlewares.GetUserID(e)

	page := 1
	if pageStr := e.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 10
	if pageSizeStr := e.QueryParam("pageSize"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	requestQueryParams := queryparams.RequestQueryParam{
		Search:       search,
		Page:         page,
		PageSize:     pageSize,
		ProjectKey:   projectKey,
		Status:       status,
		SprintID:     sprintId,
		WorkflowType: workflowType,
	}

	requests, err := c.requestService.FindAllRequestHandler(ctx, requestQueryParams, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return e.JSON(http.StatusOK, responses.Response{
		Message: "Success",
		Data:    requests,
	})
}

// GetRequestCount godoc
// @Summary      Get request counts
// @Description  Retrieves counts of requests grouped by status or other criteria for the logged-in user.
// @Tags         Requests
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} responses.RequestOverviewResponse // Assuming responses.RequestOverviewResponse exists
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/count [get]
func (c *RequestController) GetRequestCount(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	requestOverviewResponse, err := c.requestService.GetRequestCountHandler(ctx, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, requestOverviewResponse)
}

// GetRequestDetail godoc
// @Summary      Get request details
// @Description  Retrieves detailed information for a specific request ID for the logged-in user.
// @Tags         Requests
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "Request ID"
// @Success      200 {object} responses.RequestDetailResponse // Assuming responses.RequestDetailResponse exists
// @Failure      400 {object} string "Error message for invalid request ID"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/{id} [get]
func (c *RequestController) GetRequestDetail(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	requestId := e.Param("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid request ID: %s", requestId))
	}

	requestDetailResponse, err := c.requestService.GetRequestDetailHandler(ctx, userId, int32(requestIdInt))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, requestDetailResponse)

}

// GetRequestTasks godoc
// @Summary      Get tasks for a specific request
// @Description  Retrieves a paginated list of tasks associated with a given request ID.
// @Tags         Requests
// @Produce      json
// @Param        id path int true "Request ID"
// @Param        page query int false "Page number for pagination" default(1)
// @Param        pageSize query int false "Number of items per page" default(10)
// @Success      200 {object} responses.Response{data=responses.RequestTaskResponse} // Assuming appropriate response structure
// @Failure      400 {object} string "Error message for invalid request ID"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/{id}/tasks [get]
func (c *RequestController) GetRequestTasks(e echo.Context) error {
	ctx := e.Request().Context()

	// userId, _ := middlewares.GetUserID(e)

	requestId := e.Param("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid request ID: %s", requestId))

	}

	page := 1
	if pageStr := e.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 10
	if pageSizeStr := e.QueryParam("pageSize"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	requestTaskQueryParam := queryparams.RequestTaskQueryParam{
		Page:     page,
		PageSize: pageSize,
	}

	requestTasksResponse, err := c.requestService.GetRequestTasksHandler(ctx, int32(requestIdInt), requestTaskQueryParam)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, responses.Response{
		Message: "Success",
		Data:    requestTasksResponse,
	})
}

// GetRequestTasksByProject godoc
// @Summary      Get request tasks by project and filters
// @Description  Retrieves a paginated list of tasks across requests, filtered by project, type, status, and workflow type for the logged-in user.
// @Tags         Requests
// @Produce      json
// @Security     ApiKeyAuth
// @Param        projectKey query string true "Filter by project key" // Marked as required based on context
// @Param        workflowType query string false "Filter by workflow type"
// @Param        status query string false "Filter by task status"
// @Param        type query string false "Filter by task type"
// @Param        page query int false "Page number for pagination" default(1)
// @Param        pageSize query int false "Number of items per page" default(10)
// @Success      200 {object} responses.Response{data=responses.RequestTaskResponse} // Assuming appropriate response structure
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/tasks/by-project [get] // Consider a more descriptive path like /projects/{projectKey}/tasks
func (c *RequestController) GetRequestTasksByProject(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	page := 1
	if pageStr := e.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 10
	if pageSizeStr := e.QueryParam("pageSize"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	workflowType := e.QueryParam("workflowType")
	status := e.QueryParam("status")
	typeQuery := e.QueryParam("type")
	projectKey := e.QueryParam("projectKey")

	requestTaskProjectQueryParam := queryparams.RequestTaskProjectQueryParam{
		Page:         page,
		PageSize:     pageSize,
		WorkflowType: workflowType,
		Status:       status,
		Type:         typeQuery,
		ProjectKey:   projectKey,
	}

	requestTasksResponse, err := c.requestService.GetRequestTasksByProjectHandler(ctx, requestTaskProjectQueryParam, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, responses.Response{
		Message: "Success",
		Data:    requestTasksResponse,
	})
}

// GetRequestOverview godoc
// @Summary      Get request overview
// @Description  Retrieves overview information (like counts, statuses) for a specific request ID.
// @Tags         Requests
// @Produce      json
// @Param        id path int true "Request ID"
// @Success      200 {object} responses.RequestOverviewResponse // Assuming responses.RequestOverviewData exists
// @Failure      400 {object} string "Error message for invalid request ID"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/{id}/overview [get]
func (c *RequestController) GetRequestOverview(e echo.Context) error {
	ctx := e.Request().Context()

	requestId := e.Param("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid request ID: %s", requestId))
	}

	requestOverviewResponse, err := c.requestService.GetRequestOverviewHandler(ctx, int32(requestIdInt))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, requestOverviewResponse)
}

// FindAllSubRequestByRequestId godoc
// @Summary      Find all sub-requests for a given request
// @Description  Retrieves a paginated list of sub-requests associated with a specific parent request ID.
// @Tags         Requests
// @Produce      json
// @Param        id path int true "Parent Request ID"
// @Param        page query int false "Page number for pagination" default(1)
// @Param        pageSize query int false "Number of items per page" default(10)
// @Success      200 {object} responses.Response{data=responses.RequestSubRequest} // Assuming appropriate response structure
// @Failure      400 {object} string "Error message for invalid request ID"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/{id}/sub-requests [get]
func (c *RequestController) FindAllSubRequestByRequestId(e echo.Context) error {
	ctx := e.Request().Context()

	requestId := e.Param("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid request ID: %s", requestId))
	}

	page := 1
	if pageStr := e.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 10
	if pageSizeStr := e.QueryParam("pageSize"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	requestSubRequestQueryParam := queryparams.RequestSubRequestQueryParam{
		Page:     page,
		PageSize: pageSize,
	}

	subRequests, err := c.requestService.FindAllSubRequestByRequestId(ctx, int32(requestIdInt), requestSubRequestQueryParam)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, responses.Response{
		Message: "Success",
		Data:    subRequests,
	})
}

// UpdateRequest godoc
// @Summary      Update a request
// @Description  Updates details of an existing request identified by its ID. Requires user authentication.
// @Tags         Requests
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "Request ID"
// @Param        request body requests.RequestUpdateRequest true "Request Update Data"
// @Success      200 {object} string "Success message: Request updated successfully: {id}"
// @Failure      400 {object} string "Error message for invalid request ID or bad request body"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/{id} [put] // Or PATCH
func (c *RequestController) UpdateRequest(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	requestId := e.Param("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid request ID: %s", requestId))
	}

	req := new(requests.RequestUpdateRequest)
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = c.requestService.UpdateRequestHandler(ctx, int32(requestIdInt), req, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, fmt.Sprintf("Request updated successfully: %d", requestIdInt))
}

// GetRequestTasksCount godoc
// @Summary      Get count of request tasks based on filters
// @Description  Retrieves the total count of tasks matching filter criteria like project, workflow type, and task type for the logged-in user.
// @Tags         Requests
// @Produce      json
// @Security     ApiKeyAuth
// @Param        projectKey query string true "Filter by project key" // Marked as required
// @Param        workflowType query string false "Filter by workflow type"
// @Param        type query string false "Filter by task type"
// @Success      200 {object} responses.RequestTaskCountResponse // Assuming responses.RequestTaskCountResponse exists
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/tasks/count [get] // Consider a path like /projects/{projectKey}/tasks/count
func (c *RequestController) GetRequestTasksCount(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	requestTaskCount := queryparams.RequestTaskCount{
		WorkflowType: e.QueryParam("workflowType"),
		Type:         e.QueryParam("type"),
		ProjectKey:   e.QueryParam("projectKey"),
	}

	requestTaskCountResponse, err := c.requestService.GetRequestTaskCount(ctx, userId, requestTaskCount)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, requestTaskCountResponse)
}

// GetRequestCompletedForm godoc
// @Summary      Get completed form for a request
// @Description  Retrieves the completed form for a specific request ID.
// @Tags         Requests
// @Produce      json
// @Param        id path int true "Request ID"
// @Success      200 {object} responses.RequestCompletedFormResponse // Assuming responses.RequestCompletedFormResponse exists
// @Failure      400 {object} string "Error message for invalid request ID"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/{id}/completed-form [get]
func (c *RequestController) GetRequestCompletedForm(e echo.Context) error {
	ctx := e.Request().Context()

	requestId := e.Param("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid request ID: %s", requestId))
	}

	page := 1
	if pageStr := e.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 10
	if pageSizeStr := e.QueryParam("pageSize"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	requestCompletedFormQueryParam := queryparams.RequestTaskQueryParam{
		Page:     page,
		PageSize: pageSize,
	}

	requestCompletedFormResponse, err := c.requestService.GetRequestCompletedFormHandler(ctx, int32(requestIdInt), requestCompletedFormQueryParam)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, responses.Response{
		Message: "Success",
		Data:    requestCompletedFormResponse,
	})
}

// GetRequestFileManager godoc
// @Summary      Get file manager for a request
// @Description  Retrieves the file manager for a specific request ID.
// @Tags         Requests
// @Produce      json
// @Param        id path int true "Request ID"
// @Success      200 {object} responses.RequestFileManagerResponse // Assuming responses.RequestFileManagerResponse exists
// @Failure      400 {object} string "Error message for invalid request ID"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/{id}/file-manager [get]
func (c *RequestController) GetRequestFileManager(e echo.Context) error {
	ctx := e.Request().Context()

	requestId := e.Param("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid request ID: %s", requestId))
	}

	page := 1
	if pageStr := e.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 10
	if pageSizeStr := e.QueryParam("pageSize"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	requestFileManagerQueryParam := queryparams.RequestTaskQueryParam{
		Page:     page,
		PageSize: pageSize,
	}

	requestFileManagerResponse, err := c.requestService.GetRequestFileManagerHandler(ctx, int32(requestIdInt), requestFileManagerQueryParam)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, responses.Response{
		Message: "Success",
		Data:    requestFileManagerResponse,
	})
}

// GetRequestCompletedFormApproval godoc
// @Summary      Get completed form approval for a request
// @Description  Retrieves the completed form approval for a specific request ID.
// @Tags         Requests
// @Produce      json
// @Param        id path int true "Request ID"
// @Success      200 {object} responses.RequestCompletedFormApprovalResponse // Assuming responses.RequestCompletedFormApprovalResponse exists
// @Failure      400 {object} string "Error message for invalid request ID"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/{id}/completed-form/approval [get]
func (c *RequestController) GetRequestCompletedFormApproval(e echo.Context) error {
	ctx := e.Request().Context()

	requestId := e.Param("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid request ID: %s", requestId))
	}

	dataId := e.Param("dataId")
	if dataId == "" {
		return e.JSON(http.StatusBadRequest, "Data ID is required")
	}

	requestCompletedFormApprovalResponse, err := c.requestService.GetRequestCompletedFormApprovalHandler(ctx, int32(requestIdInt), dataId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, requestCompletedFormApprovalResponse)
}

// FindAllHistoryByRequestId godoc
// @Summary      Find all history by request id
// @Description  Find all history by request id
// @Tags         History
// @Accept       json
// @Produce      json
// @Param        requestId path int true "Request ID"
// @Success      200 {object} []responses.HistoryResponse "History"
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /history/{requestId} [get]
func (c *RequestController) FindAllHistoryByRequestId(e echo.Context) error {
	ctx := e.Request().Context()

	requestId := e.Param("requestId")

	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	history, err := c.requestService.FindAllHistoryByRequestId(ctx, int32(requestIdInt))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, history)
}

// ReportMidSprintTasks godoc
// @Summary      Report mid sprint tasks
// @Description  Report mid sprint tasks
// @Tags         Requests
// @Produce      json
// @Success      200 {object} []responses.RequestTaskResponse "Request Task Response"
// @Failure      500 {object} string "Error message for internal server error"
// @Router       /requests/report/mid-sprint-tasks [get]
func (c *RequestController) ReportMidSprintTasks(e echo.Context) error {
	ctx := e.Request().Context()

	startTime := e.QueryParam("startTime")
	endTime := e.QueryParam("endTime")

	if startTime == "" || endTime == "" {
		return e.JSON(http.StatusBadRequest, "Start time and end time are required")
	}

	startTimeInt, err := strconv.Atoi(startTime)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	endTimeInt, err := strconv.Atoi(endTime)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	requestTasksResponse, err := c.requestService.ReportMidSprintTasks(ctx, time.Unix(int64(startTimeInt), 0), time.Unix(int64(endTimeInt), 0))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, requestTasksResponse)
}
