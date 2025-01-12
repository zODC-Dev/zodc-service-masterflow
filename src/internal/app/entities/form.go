package entities

import (
	"gorm.io/datatypes"
)

type Category string

const (
	BASIC_FIELD     Category = "BASIC_FIELD"
	DATE_TIME_FIELD Category = "DATE_TIME_FIELD"
	ADVANCED_FIELD  Category = "ADVANCED_FIELD"
)

type Form struct {
	BaseModel

	FileId      string         `json:"fileId" gorm:"not null"`
	FileName    string         `json:"fileName" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	Function    string         `json:"function" gorm:"not null"`
	Template    string         `json:"template" gorm:"not null"`
	DataSheet   datatypes.JSON `json:"dataSheet" gorm:"type:jsonb;not null"`
	Description string         `json:"description" gorm:"not null"`
	FormFields  []FormField    `json:"formFields" gorm:"foreignKey:FormID;not null"`
}

type FormField struct {
	BaseModel

	Icon            string         `json:"icon" gorm:"not null"`
	Title           string         `json:"title" gorm:"not null"`
	Category        Category       `json:"category" gorm:"not null"`
	FieldName       string         `json:"fieldName" gorm:"not null"`
	FieldType       string         `json:"fieldType" gorm:"not null"`
	Required        bool           `json:"required" gorm:"not null"`
	AdvancedOptions datatypes.JSON `json:"advancedOptions" gorm:"type:jsonb;not null"`
	FormID          uint           `json:"formId" gorm:"not null"`
	ColNum          uint           `json:"colNum" gorm:"not null"`
}
