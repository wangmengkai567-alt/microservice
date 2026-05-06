package model

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"not null"`
	Stock       int       `gorm:"not null;default:0"`
	ImageUrl    string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
