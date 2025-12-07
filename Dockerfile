# ======================================================
# 1) 构建阶段：使用镜像加速源的 golang alpine 构建二进制
#    使用腾讯云公共镜像加速，规避 dockerhub 访问超时
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

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o main .

# ======================================================
# 2) 运行阶段：使用 scratch，避免拉取运行时基础镜像
#    从构建阶段拷贝证书与可执行文件
# ======================================================
FROM alpine:3.13

WORKDIR /app


COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/main /app/main
COPY index.html ./index.html

# 以非特权用户运行（直接设置 uid），避免 80 端口权限问题，应用监听 8080
USER 10001

EXPOSE 8080

ENTRYPOINT ["/app/main"]
