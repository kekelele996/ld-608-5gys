FROM golang:1.22-alpine AS build
WORKDIR /app
# 容器内使用纯 Go MySQL 驱动，CGO_ENABLED=0 产出静态二进制。
ENV CGO_ENABLED=0 GOOS=linux
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./main.go

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache wget ca-certificates
COPY --from=build /app/server ./server
CMD ["./server"]
