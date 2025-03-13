package users

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey"`
	Name       string         `gorm:"column:name;not null"`
	Username   string         `gorm:"column:username;uniqueIndex;not null"`
	Email      string         `gorm:"column:email;uniqueIndex;not null"`
	Password   string         `gorm:"column:password;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	ModifiedAt time.Time      `gorm:"column:modified_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
