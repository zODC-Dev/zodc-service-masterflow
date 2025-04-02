package responses

import "time"

type FormTemplateFieldsFindAll struct {
	ID              int32                  `json:"id"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	DeletedAt       *time.Time             `json:"deletedAt"`
	FieldID         string                 `json:"fieldId"`
	Icon            string                 `json:"icon"`
	Title           string                 `json:"title"`
	Category        string                 `json:"category"`
	FieldName       string                 `json:"fieldName"`
	FieldType       string                 `json:"fieldType"`
	Required        bool                   `json:"required"`
	AdvancedOptions map[string]interface{} `json:"advancedOptions"`
	ColNum          int32                  `json:"colNum"`
	FormID          int32                  `json:"formId"`
}

type FormTemplateFindAll struct {
	ID        int32           `json:"id"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt *time.Time      `json:"deletedAt"`
	FileName  string          `json:"fileName"`
	Title     string          `json:"title"`
	Category  CategoryFindAll `json:"category"`

	Version int32 `json:"version"`

	TemplateID  *int32                  `json:"templateId"`
	DataSheet   *map[string]interface{} `json:"dataSheet"`
	Description string                  `json:"description"`
	Decoration  string                  `json:"decoration"`
}
