package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subscription represents the database model for user server
type Subscription struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ServiceName string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_service_unique;index"`
	Price       int       `gorm:"not null;check:price > 0"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_service_unique"`

	// Date fields
	StartMonth int  `gorm:"not null;check:start_month >= 1 AND start_month <= 12;index:idx_start_date"`
	StartYear  int  `gorm:"not null;index:idx_start_date"`
	EndMonth   *int `gorm:"check:end_month >= 1 AND end_month <= 12;index:idx_end_date"`
	EndYear    *int `gorm:"index:idx_end_date"`
}

// TableName specifies the table name
func (*Subscription) TableName() string {
	return "subscriptions"
}

// BeforeCreate GORM hook
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
