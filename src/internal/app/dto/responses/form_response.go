package responses

import (
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
	"gorm.io/datatypes"
)

type FormResponse struct {
	entities.BaseModel

	FileId      string                 `json:"fileId"`
	FileName    string                 `json:"fileName"`
	Title       string                 `json:"title"`
	Function    string                 `json:"function"`
	Template    string                 `json:"template"`
	DataSheet   datatypes.JSON         `json:"dataSheet"`
	Description string                 `json:"description"`
	FormFields  [][]entities.FormField `json:"formFields"`
}
