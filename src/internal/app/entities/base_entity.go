package entities

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey;autoIncrement;unique"`
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"deleteAt" gorm:"index"`
}
