
FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .


RUN go build -o jwt-mongo-auth-service cmd/main.go


FROM alpine:latest


RUN apk add --no-cache ca-certificates && update-ca-certificates && \
    apk add --no-cache libc6-compat


COPY --from=build /app/jwt-mongo-auth-service /jwt-mongo-auth-service


EXPOSE 8080


CMD ["/jwt-mongo-auth-service"]
