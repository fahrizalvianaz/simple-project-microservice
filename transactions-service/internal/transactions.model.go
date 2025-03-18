package internal

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID         uint           `gorm:"primaryKey"`
	BookID     uint           `gorm:"column:book_id;not null"`
	UserID     uint           `gorm:"column:user_id;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	ModifiedAt time.Time      `gorm:"column:modified_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Book) TableName() string {
	return "transactions"
}
