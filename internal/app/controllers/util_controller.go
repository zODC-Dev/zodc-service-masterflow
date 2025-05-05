package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

// ExcelExtract godoc
// @Summary      Extract data from an Excel file
// @Description  Uploads an Excel file and extracts data from it.
// @Tags         Utils
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Excel file to upload"
// @Success      200 {object} map[string]interface{} // Or a more specific DTO if the structure is known
// @Failure      400 {object} map[string]string "error: No file is received"
// @Failure      500 {object} map[string]string "error: Fail extract excel file"
// @Router       /utils/excel-extract [post]
func ExcelExtract(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "No file is received"})
	}

	result, err := utils.ExcelExtract(file)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Fail extract excel file"})
	}

	return ctx.JSON(http.StatusOK, result)
}
