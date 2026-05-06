package handler

import (
	"api-gateway/internal/client"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pb "api-gateway/proto/product"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// UploadImage 上传商品图片
func (h *ProductHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported image format"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image size must be less than 5MB"})
		return
	}

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, fileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	imageURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, c.Request.Host, fileName)

	c.JSON(http.StatusOK, gin.H{
		"image_url": imageURL,
	})
}

// CreateProduct 创建商品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" binding:"required,gte=0"`
		ImageUrl    string  `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取 token
	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	// 创建带 token 的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// 调用 Product Service gRPC
	productClient := pb.NewProductServiceClient(client.Clients.ProductConn)
	resp, err := productClient.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int32(req.Stock),
		ImageUrl:    req.ImageUrl,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   resp.Id,
		"name": req.Name,
	})
}

// GetProduct 获取商品详情
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用 Product Service gRPC
	productClient := pb.NewProductServiceClient(client.Clients.ProductConn)
	resp, err := productClient.GetProduct(ctx, &pb.GetProductRequest{
		Id: uint32(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          resp.Id,
		"name":        resp.Name,
		"description": resp.Description,
		"price":       resp.Price,
		"stock":       resp.Stock,
		"image_url":   resp.ImageUrl,
	})
}

// UpdateProduct 更新商品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		ImageUrl    string  `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	productClient := pb.NewProductServiceClient(client.Clients.ProductConn)
	_, err = productClient.UpdateProduct(ctx, &pb.UpdateProductRequest{
		Id:          uint32(id),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int32(req.Stock),
		ImageUrl:    req.ImageUrl,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product updated successfully",
	})
}

// DeleteProduct 删除商品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	token, _ := c.Get("token")
	tokenStr, _ := token.(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tokenStr != "" {
		md := metadata.Pairs("authorization", tokenStr)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	productClient := pb.NewProductServiceClient(client.Clients.ProductConn)
	_, err = productClient.DeleteProduct(ctx, &pb.DeleteProductRequest{
		Id: uint32(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product deleted successfully",
	})
}

// ListProducts 商品列表
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	productClient := pb.NewProductServiceClient(client.Clients.ProductConn)
	resp, err := productClient.ListProducts(ctx, &pb.ListProductsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type Product struct {
		ID          uint32  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int32   `json:"stock"`
		ImageUrl    string  `json:"image_url"`
	}

	products := make([]Product, 0, len(resp.Products))
	for _, p := range resp.Products {
		products = append(products, Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			ImageUrl:    p.ImageUrl,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    resp.Total,
	})
}
