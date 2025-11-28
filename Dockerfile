# ======================================================
# 1) 构建阶段：使用 Go 官方镜像构建二进制
# ======================================================
FROM golang:1.22-alpine AS builder

# 使用国内代理提高构建速度（可选）
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 先复制 go.mod / go.sum，加速缓存
COPY go.mod go.sum ./
RUN go mod download

# 再复制源代码
COPY . .

# 编译为 Linux 可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ======================================================
# 2) 运行阶段：使用轻量级 Alpine 镜像
# ======================================================
FROM alpine:3.13

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
