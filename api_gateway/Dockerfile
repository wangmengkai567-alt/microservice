# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 设置 Go 代理（多级回退，避免单点失败）
ENV GOPROXY=https://proxy.golang.org,https://goproxy.cn,direct

RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o api-gateway cmd/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/api-gateway .
COPY --from=builder /app/config ./config

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./api-gateway"]
