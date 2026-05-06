package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TraceInterceptor 添加请求追踪ID
func TraceInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// 从metadata中获取trace_id，如果没有则生成新的
		traceID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids := md.Get("trace_id"); len(ids) > 0 {
				traceID = ids[0]
			}
		}

		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 将trace_id添加到context中
		ctx = context.WithValue(ctx, "trace_id", traceID)

		return handler(ctx, req)
	}
}

// GetTraceID 从context中获取trace_id
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}
