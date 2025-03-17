package dto

type CreateRequest struct {
	Title       string `json:"title" binding:"required" example:"Rich Dad, Poor Dad"`
	Author      string `json:"author" binding:"required" example:"Robert T Kiyosaki"`
	Description string `json:"description" binding:"required" example:"Learn from rich dad and poor dad about financial management"`
	Price       int    `json:"price" binding:"required" example:"100000"`
	Stock       int    `json:"stock" binding:"required" example:"50"`
}
