package entities

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time      `json:"UpdatedAt" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"deleteAt" gorm:"index"`
}
