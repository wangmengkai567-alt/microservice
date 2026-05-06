package middleware

import (
	"context"
	"strings"

	"product-service/config"
	"product-service/internal/pkg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// GetProduct、ListProducts 和 AdjustStock 接口不需要认证
		// GetProduct 和 ListProducts 是只读操作
		// AdjustStock 是服务间调用（订单服务调用）
		if strings.Contains(info.FullMethod, "GetProduct") ||
			strings.Contains(info.FullMethod, "ListProducts") ||
			strings.Contains(info.FullMethod, "AdjustStock") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokenArr := md["authorization"]
		if len(tokenArr) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		token := tokenArr[0]

		claims, err := pkg.ParseToken(token, config.Global.JWT.Secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		newCtx := context.WithValue(ctx, "userID", claims.UserID)

		return handler(newCtx, req)
	}
}
