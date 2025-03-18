package api

import (
	"book-service/internal"
	"book-service/internal/api/dto"
	"net/http"
	"strconv"

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

func (b *BookHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		genericResponse.BadRequestResponse(ctx, "Invalid ID format", err.Error())
		return
	}
	response, err := b.bookService.FindByID(ctx.Request.Context(), uint(idUint))
	if err != nil {
		genericResponse.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	genericResponse.SuccessResponse(ctx, 200, "Book found", response)

}
