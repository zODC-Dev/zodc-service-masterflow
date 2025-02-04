package requests

type FormFieldsCreate struct {
	FieldID         string
	Icon            string
	Title           string
	Category        string
	FieldName       string
	FieldType       string
	Required        bool
	AdvancedOptions map[string]interface{}
	ColNum          int32
	FormID          int32
}

type FormCreate struct {
	FileName    string
	Title       string
	CategoryID  *int32
	Version     int32
	TemplateID  *int32
	DataSheet   *map[string]interface{}
	Description string
	Decoration  string
	FormFields  [][]FormFieldsCreate `json:"formFields"`
}
