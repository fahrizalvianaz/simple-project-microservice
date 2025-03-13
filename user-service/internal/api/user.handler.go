package api

import (
	users "bookstore-framework/internal"
	"bookstore-framework/internal/api/dto"
	"net/http"

	genericResponse "github.com/fahrizalvianaz/shared-response/httputil"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService users.UserService
}

func NewUserHandler(userService users.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterHandler godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body     dto.RegisterRequest true "User information"
// @Success      201  {object}    pkg.Response{data=dto.RegisterResponse} "User registered successfully"
// @Failure      400  {object}    pkg.Response "Invalid Request format"
// @Router       /users/register [post]
func (h *UserHandler) RegisterHandler(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		genericResponse.BadRequestResponse(ctx, "Invalid Request format", err.Error())
		return
	}

	response, err := h.userService.Register(ctx.Request.Context(), req)
	if err != nil {
		genericResponse.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	genericResponse.CreatedResponse(ctx, "User registered successfully", response)
}

// LoginHandler godoc
// @Summary      Login user
// @Description  Login user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body     dto.LoginRequest true "User information"
// @Success      201  {object}    pkg.Response{data=dto.LoginResponse} "Login successfully"
// @Failure      400  {object}    pkg.Response "Invalid Request format"
// @Router       /users/login [post]
func (h *UserHandler) LoginHandler(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		genericResponse.BadRequestResponse(ctx, "Invalid Request format", err.Error())
		return
	}

	response, err := h.userService.Login(ctx.Request.Context(), req)
	if err != nil {
		genericResponse.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	genericResponse.OkResponse(ctx, "Login Successfully", response)
}

// LoginHandler godoc
// @Summary      Get User
// @Description  Get user account
// @Tags         users
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Success      201  {object}    pkg.Response{data=dto.ProfileResponse} "Profile retrieve successfully"
// @Failure      400  {object}    pkg.Response "Invalid Request format"
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userID, exist := ctx.Get("userID")
	if !exist {
		genericResponse.ErrorResponse(ctx, http.StatusUnauthorized, "User not found", nil)
		return
	}

	profile, err := h.userService.GetProfile(ctx, userID.(uint))
	if err != nil {
		genericResponse.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve profile", err.Error())
		return
	}

	genericResponse.OkResponse(ctx, "Profile retrieve successfully", profile)
}
