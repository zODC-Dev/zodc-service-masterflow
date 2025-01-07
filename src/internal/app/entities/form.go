package entities

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Category string

const (
	BASIC_FIELD     Category = "BASIC_FIELD"
	DATE_TIME_FIELD Category = "DATE_TIME_FIELD"
	ADVANCED_FIELD  Category = "ADVANCED_FIELD"
)

type Form struct {
	gorm.Model

	Icon            string         `json:"icon"`
	Title           string         `json:"title"`
	Category        Category       `json:"category"`
	FieldName       string         `json:"fieldName"`
	FieldType       string         `json:"fieldType"`
	Required        bool           `json:"required"`
	AdvancedOptions datatypes.JSON `json:"advancedOptions" gorm:"type:jsonb"`
	FormExcelID     uint           `json:"formExcelId" gorm:"index"`
}

type FormExcel struct {
	gorm.Model

	FileName    string `json:"fileName"`
	Title       string `json:"title"`
	Function    string `json:"function"`
	Template    string `json:"template"`
	DataSheet   string `json:"dataSheet"`
	Description string `json:"description"`
	Forms       []Form `json:"formDetails" gorm:"foreignKey:FormExcelID"`
}
