package dto

// RegisterRequest represents a registration request
// @Description Registration request payload
type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"johndoe"`
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required" example:"johndoe@gmail.com"`
	Password string `json:"password" binding:"required" example:"xxxxxxx"`
}

// LoginRequest represents a registration request
// @Description Login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"xxxxxxx"`
}
