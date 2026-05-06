package model

import "time"

type Order struct {
	ID          uint      `gorm:"primaryKey"`
	OrderNo     string    `gorm:"unique;not null;index"`
	UserID      uint      `gorm:"not null;index"`
	ProductID   uint      `gorm:"not null"`
	ProductName string    `gorm:"not null"`
	Quantity    int       `gorm:"not null"`
	TotalPrice  float64   `gorm:"not null"`
	Status      string    `gorm:"not null;default:'pending'"` // pending, paid, completed, cancelled
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
