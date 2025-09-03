package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subscription represents the database model for user server
type Subscription struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Required fields from TZ
	ServiceName string    `gorm:"type:varchar(255);not null"` // Просто строка
	Price       int       `gorm:"not null;check:price > 0"`   // Целое число рублей
	UserID      uuid.UUID `gorm:"type:uuid;not null"`         // Просто храним UUID

	// Date fields - храним отдельно месяц и год для запросов
	StartMonth int  `gorm:"not null;check:start_month >= 1 AND start_month <= 12"`
	StartYear  int  `gorm:"not null;check:start_year >= 2020"`
	EndMonth   *int `gorm:"check:end_month >= 1 AND end_month <= 12"`
	EndYear    *int `gorm:"check:end_year >= 2020"`

	// Составной индекс для частых запросов
	IndexUserService *struct{} `gorm:"uniqueIndex:idx_user_service"` // Уникальная подписка на сервис
	IndexUserID      *struct{} `gorm:"index:idx_user_id"`
	IndexServiceName *struct{} `gorm:"index:idx_service_name"`
	IndexDate        *struct{} `gorm:"index:idx_date"` // Для фильтрации по датам
}

// TableName specifies the table name
func (*Subscription) TableName() string {
	return "server"
}

// BeforeCreate GORM hook
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
