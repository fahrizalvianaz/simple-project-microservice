package internal

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID          uint           `gorm:"primaryKey"`
	Title       string         `gorm:"column:title;not null"`
	Author      string         `gorm:"column:author;uniqueIndex;not null"`
	Description string         `gorm:"column:description;uniqueIndex;not null"`
	Price       int            `gorm:"column:price;not null"`
	Stock       int            `gorm:"column:stock;not null"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	ModifiedAt  time.Time      `gorm:"column:modified_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Book) TableName() string {
	return "books"
}
