package controllers

import (
	"fmt"
	"net/http"
	"strconv"

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

func (c *RequestController) GetRequestCount(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	requestOverviewResponse, err := c.requestService.GetRequestCountHandler(ctx, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, requestOverviewResponse)
}

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
