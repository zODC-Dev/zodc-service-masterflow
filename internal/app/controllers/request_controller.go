package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
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
		Search:     search,
		Page:       page,
		PageSize:   pageSize,
		ProjectKey: projectKey,
		Status:     status,
	}

	requests, err := c.requestService.FindAllRequestHandler(ctx, requestQueryParams, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return e.JSON(http.StatusOK, requests)
}

func (c *RequestController) GetRequestOverview(e echo.Context) error {
	ctx := e.Request().Context()

	userId, _ := middlewares.GetUserID(e)

	requestOverviewResponse, err := c.requestService.GetRequestOverviewHandler(ctx, userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, requestOverviewResponse)
}
