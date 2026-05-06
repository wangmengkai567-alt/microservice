package handler

import (
	"context"
	"product-service/internal/logger"
	"product-service/internal/middleware"
	"product-service/internal/service"
	pb "product-service/proto"

	"go.uber.org/zap"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	traceID := middleware.GetTraceID(ctx)
	
	logger.Info("Creating product",
		zap.String("trace_id", traceID),
		zap.String("name", req.Name),
		zap.Float64("price", req.Price),
		zap.Int32("stock", req.Stock),
	)

	id, err := h.service.CreateProduct(req.Name, req.Description, req.Price, int(req.Stock), req.ImageUrl)
	if err != nil {
		logger.Error("Failed to create product",
			zap.String("trace_id", traceID),
			zap.String("name", req.Name),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Product created successfully",
		zap.String("trace_id", traceID),
		zap.Uint("product_id", id),
	)

	return &pb.CreateProductResponse{
		Id:      uint32(id),
		Message: "product created successfully",
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := h.service.GetProduct(uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{
		Id:          uint32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       int32(product.Stock),
		ImageUrl:    product.ImageUrl,
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	traceID := middleware.GetTraceID(ctx)
	
	logger.Info("Updating product",
		zap.String("trace_id", traceID),
		zap.Uint32("product_id", req.Id),
	)

	err := h.service.UpdateProduct(uint(req.Id), req.Name, req.Description, req.Price, int(req.Stock), req.ImageUrl)
	if err != nil {
		logger.Error("Failed to update product",
			zap.String("trace_id", traceID),
			zap.Uint32("product_id", req.Id),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Product updated successfully",
		zap.String("trace_id", traceID),
		zap.Uint32("product_id", req.Id),
	)

	return &pb.UpdateProductResponse{
		Message: "product updated successfully",
	}, nil
}

func (h *ProductHandler) AdjustStock(ctx context.Context, req *pb.AdjustStockRequest) (*pb.AdjustStockResponse, error) {
	traceID := middleware.GetTraceID(ctx)
	
	logger.Info("Adjusting stock",
		zap.String("trace_id", traceID),
		zap.Uint32("product_id", req.Id),
		zap.Int32("delta", req.Delta),
	)

	err := h.service.AdjustStock(uint(req.Id), int(req.Delta))
	if err != nil {
		logger.Error("Failed to adjust stock",
			zap.String("trace_id", traceID),
			zap.Uint32("product_id", req.Id),
			zap.Int32("delta", req.Delta),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Stock adjusted successfully",
		zap.String("trace_id", traceID),
		zap.Uint32("product_id", req.Id),
		zap.Int32("delta", req.Delta),
	)

	return &pb.AdjustStockResponse{
		Message: "stock adjusted successfully",
	}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	traceID := middleware.GetTraceID(ctx)
	
	logger.Info("Deleting product",
		zap.String("trace_id", traceID),
		zap.Uint32("product_id", req.Id),
	)

	err := h.service.DeleteProduct(uint(req.Id))
	if err != nil {
		logger.Error("Failed to delete product",
			zap.String("trace_id", traceID),
			zap.Uint32("product_id", req.Id),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Product deleted successfully",
		zap.String("trace_id", traceID),
		zap.Uint32("product_id", req.Id),
	)

	return &pb.DeleteProductResponse{
		Message: "product deleted successfully",
	}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := h.service.ListProducts(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          uint32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       int32(p.Stock),
			ImageUrl:    p.ImageUrl,
		})
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    int32(total),
	}, nil
}
