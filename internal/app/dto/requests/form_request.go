package requests

import "encoding/json"

type FormTemplateFieldsCreate struct {
	FieldID         string
	Icon            string
	Title           string
	Category        string
	FieldName       string
	FieldType       string
	Required        bool
	Readonly        bool
	AdvancedOptions *json.RawMessage
	ColNum          int32
	FormID          int32
}

type FormTemplateCreate struct {
	FileName    string
	Title       string
	CategoryID  *int32
	TemplateID  *int32
	DataSheet   *json.RawMessage
	Description string
	Decoration  string
	FormFields  [][]FormTemplateFieldsCreate `json:"formFields"`
}

type FormTemplateUpdate struct {
	FileName    string
	Title       string
	CategoryID  *int32
	TemplateID  *int32
	DataSheet   *json.RawMessage
	Description string
	Decoration  string
}
