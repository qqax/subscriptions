package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	StartDate time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	EndDate   *time.Time     `gorm:"default:null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Service   Service
	ID        uint      `gorm:"primaryKey"`
	ServiceID uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;unique"`
}
