package middleware

import (
	"context"

	"order-service/config"
	"order-service/internal/pkg"

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

		// 所有订单接口都需要认证
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

		// 将 userID 存入 context
		newCtx := context.WithValue(ctx, "userID", claims.UserID)

		return handler(newCtx, req)
	}
}
