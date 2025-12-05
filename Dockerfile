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
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o main .

# ======================================================
# 2) 运行阶段：使用轻量级 Alpine 镜像
# ======================================================
FROM alpine:3.18

WORKDIR /app

# 基础运行环境（证书、时区），保证 HTTPS 能正常访问
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

COPY --from=builder /app/main .
# 业务静态资源（主页）
COPY index.html ./index.html

# 以非 root 运行，提升安全性
RUN adduser -D -u 10001 appuser
USER appuser

# 暴露服务端口
EXPOSE 80

CMD ["./main"]
