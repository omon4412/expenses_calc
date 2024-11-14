package models

import (
	"time"
)

type Expense struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"expense_id"`
	Name       string    `gorm:"not null" json:"name"`
	UserID     uint      `gorm:"not null" json:"-"`
	User       User      `gorm:"foreignKey:UserID" json:"-"`
	CategoryID uint      `gorm:"not null" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"-"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Date       time.Time `gorm:"not null" json:"date"`
}
