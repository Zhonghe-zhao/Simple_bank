# Build stage 构建二进制文件
FROM golang:1.23-alpine3.20 AS build
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Production stage 生产环境
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]

