# build stage
FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# run stage (минимальный образ)
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/server .
ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ["/app/server"]
