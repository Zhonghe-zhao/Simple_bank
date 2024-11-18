# Build stage
FROM golang:1.23-alpine3.20 AS build
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L  https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && mv migrate /app/migrate
     

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/migrate /app/migrate
COPY app.env .
COPY start.sh . 
COPY wait-for.sh . 
COPY db/migration ./migration

EXPOSE 8080 
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]

