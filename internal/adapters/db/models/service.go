package models

import "time"

type Service struct {
	ServiceName   string         `gorm:"not null"`
	Subscriptions []Subscription `gorm:"foreignKey:ServiceID"`
	ID            uint           `gorm:"primaryKey"`
}

type Price struct {
	Date         time.Time `gorm:"not null"`
	Service      Service
	ID           uint    `gorm:"primaryKey"`
	ServiceID    uint    `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MonthlyPrice float64 `gorm:"not null"`
}

type ServiceStatus struct {
	Date      time.Time `gorm:"not null"`
	Service   Service
	ID        uint `gorm:"primaryKey"`
	ServiceID uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Active    bool `gorm:"not null"`
}
