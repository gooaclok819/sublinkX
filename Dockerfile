# Build stage
FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o sublinkX

# Final stage
FROM alpine:latest
WORKDIR /app

# 设置时区为 Asia/Shanghai
ENV TZ=Asia/Shanghai

COPY --from=builder /app/sublinkX /app/sublinkX
EXPOSE 8000
CMD ["/app/sublinkX"]

