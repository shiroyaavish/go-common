package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel is the basic model that is used throughout
//
// proteus:generate
type BaseModel struct {
	ID        *uuid.UUID     `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
