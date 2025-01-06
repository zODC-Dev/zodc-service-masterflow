package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/src/pkg/utils"
)

func ExcelExtract(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "No file is received"})
	}

	result, err := utils.ExcelExtract(file)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Fail extract excel file"})
	}

	return ctx.JSON(http.StatusOK, result)
}
