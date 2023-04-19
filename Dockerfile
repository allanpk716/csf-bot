FROM golang:1.18-buster AS builder

LABEL stage=gobuilder

# 开始编译
ENV CGO_ENABLED 0
ENV GO111MODULE=on
ENV GOOS linux

# 切换工作目录
WORKDIR /homelab/buildspace
# 复制文件
COPY . .
# 执行编译，-o 指定保存位置和程序编译名称
RUN cd ./cmd/csf-bot \
    && go build -ldflags="-s -w" -o /app/csf-bot

# 运行时环境
FROM alpine
ENV TZ=Asia/Shanghai \
    PERMS=true

RUN apk add --no-cache \
       bash \
       ca-certificates \
       tini \
       su-exec \
       tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo "${TZ}" > /etc/timezone \
    && rm -rf /tmp/* /var/cache/apk/*

# 主程序
COPY --from=builder /app/csf-bot /app/csf-bot
VOLUME ["/app/etc"]
WORKDIR /app
ENTRYPOINT ["/app/csf-bot"]

