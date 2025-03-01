package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

func (c *CategoryController) FindAll(e echo.Context) error {
	ctx := e.Request().Context()

	typeQueryParam := e.QueryParam("type")

	categories, err := c.categoryService.FindAll(ctx, typeQueryParam)
	if err != nil {
		return e.JSON(http.StatusBadGateway, err.Error())
	}

	return e.JSON(http.StatusOK, categories)

}
