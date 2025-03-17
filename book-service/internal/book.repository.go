package internal

import (
	"context"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(ctx context.Context, book *Book) (*Book, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{
		db: db,
	}
}

func (b *bookRepository) Create(ctx context.Context, book *Book) (*Book, error) {
	result := b.db.WithContext(ctx).Create(&book)
	if result.Error != nil {
		return nil, result.Error
	}
	return book, nil
}
