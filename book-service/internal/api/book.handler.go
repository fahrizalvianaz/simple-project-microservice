package api

import (
	"book-service/internal"
	"book-service/internal/api/dto"
	"net/http"

	genericResponse "github.com/fahrizalvianaz/shared-response/httputil"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService internal.BookService
}

func NewBookHandler(bookService internal.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (b *BookHandler) Create(ctx *gin.Context) {
	var req dto.CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		genericResponse.BadRequestResponse(ctx, "Invalid Request format", err.Error())
		return
	}

	response, err := b.bookService.Create(ctx.Request.Context(), req)
	if err != nil {
		genericResponse.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	genericResponse.CreatedResponse(ctx, "Book registered successfully", response)
}
