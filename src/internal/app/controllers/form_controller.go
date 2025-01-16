package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database/models"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/interfaces"
)

type formControllerImpl struct {
	formService interfaces.FormService
}

func NewFormController(formService interfaces.FormService) *formControllerImpl {
	return &formControllerImpl{
		formService: formService,
	}
}

func (c *formControllerImpl) Create(ctx echo.Context) error {
	ctxReq := ctx.Request().Context()

	createFormRequest := new(models.CreateFormRequest)

	if err := ctx.Bind(createFormRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.formService.Create(ctxReq, createFormRequest); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]string{
		"message": "Form created successfully",
	})
}

func (c *formControllerImpl) FindAll(ctx echo.Context) error {
	ctxReq := ctx.Request().Context()

	forms, err := c.formService.FindAll(ctxReq)
	if err != nil {
		return ctx.JSON(http.StatusBadGateway, err.Error())
	}

	return ctx.JSON(http.StatusCreated, forms)
}

// func (c *formControllerImpl) Delete(ctx echo.Context) error {
// 	id := ctx.Param("id")

// 	form, err := c.formService.FindById(id)
// 	if err != nil {
// 		return ctx.JSON(http.StatusBadGateway, err.Error())
// 	}

// 	deleteErr := c.formService.Delete(form)
// 	if deleteErr != nil {
// 		return ctx.JSON(http.StatusBadGateway, deleteErr.Error())
// 	}

// 	return ctx.JSON(http.StatusCreated, map[string]string{
// 		"message": "Form created successfully",
// 	})
// }
