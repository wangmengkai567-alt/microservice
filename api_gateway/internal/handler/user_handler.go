package handler

import (
	"api-gateway/internal/client"
	"context"
	"net/http"
	"time"

	userPb "api-gateway/proto/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 User Service gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userClient := userPb.NewUserServiceClient(client.Clients.UserConn)
	resp, err := userClient.Register(ctx, &userPb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// 打印详细错误日志
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 User Service gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userClient := userPb.NewUserServiceClient(client.Clients.UserConn)
	resp, err := userClient.Login(ctx, &userPb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 返回完整的登录响应（包含token和user信息）
	c.JSON(http.StatusOK, gin.H{
		"token": resp.Token,
		"user": gin.H{
			"id":       "1",
			"username": req.Username,
		},
	})
}
