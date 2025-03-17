package internal

import (
	"book-service/internal/api/dto"
	"context"
	"time"
)

type BookService interface {
	Create(ctx context.Context, request dto.CreateRequest) (*dto.CreateResponse, error)
}

type bookService struct {
	bookRepository BookRepository
}

func NewBookService(bookRepository BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepository,
	}
}

func (b *bookService) Create(ctx context.Context, request dto.CreateRequest) (*dto.CreateResponse, error) {
	book := &Book{
		Title:       request.Title,
		Author:      request.Author,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	result, err := b.bookRepository.Create(ctx, book)

	if err != nil {
		return nil, err
	}

	response := &dto.CreateResponse{
		ID:          result.ID,
		Title:       result.Title,
		Author:      result.Author,
		Description: result.Description,
		Price:       result.Price,
		Stock:       result.Stock,
		CreatedAt:   time.Now(),
	}

	return response, nil
}
