package models

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	database "github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database/generated"
)

type FormFieldRequest struct {
	FieldID         string          `db:"field_id" json:"fieldId"`
	Icon            string          `db:"icon" json:"icon"`
	Title           string          `db:"title" json:"title"`
	Category        string          `db:"category" json:"category"`
	FieldName       string          `db:"field_name" json:"fieldName"`
	FieldType       string          `db:"field_type" json:"fieldType"`
	Required        bool            `db:"required" json:"required"`
	AdvancedOptions json.RawMessage `db:"advanced_options" json:"advancedOptions"`
	ColNum          int32           `db:"col_num" json:"colNum"`
	FormID          int32           `db:"form_id" json:"formId"`
}

type CreateFormRequest struct {
	FileName    string               `db:"file_name" json:"fileName"`
	Title       string               `db:"title" json:"title"`
	Function    string               `db:"function" json:"function"`
	Version     int32                `db:"version" json:"version"`
	Template    pgtype.Text          `db:"template" json:"template"`
	Datasheet   json.RawMessage      `db:"datasheet" json:"datasheet"`
	Description string               `db:"description" json:"description"`
	Decoration  string               `db:"decoration" json:"decoration"`
	FormFields  [][]FormFieldRequest `json:"formFields"`
}

type FindAllFormsResponse struct {
	ID          int32                  `db:"id" json:"id"`
	CreatedAt   pgtype.Timestamp       `db:"created_at" json:"createdAt"`
	UpdatedAt   pgtype.Timestamp       `db:"updated_at" json:"updatedAt"`
	DeletedAt   pgtype.Timestamp       `db:"deleted_at" json:"deletedAt"`
	FileName    string                 `db:"file_name" json:"fileName"`
	Title       string                 `db:"title" json:"title"`
	Function    string                 `db:"function" json:"function"`
	Version     int32                  `db:"version" json:"version"`
	Template    pgtype.Text            `db:"template" json:"template"`
	Datasheet   json.RawMessage        `db:"datasheet" json:"datasheet"`
	Description string                 `db:"description" json:"description"`
	Decoration  string                 `db:"decoration" json:"decoration"`
	FormFields  [][]database.FormField `json:"formFields"`
}
