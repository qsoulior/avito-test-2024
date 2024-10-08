FROM golang:1.23.0-alpine3.20 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./main ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=build /build/main ./
COPY --from=build /build/migrations ./migrations
EXPOSE 8080
ENV POSTGRES_MIGRATIONS migrations
ENTRYPOINT ["./main"]