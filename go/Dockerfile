FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM alpine:3.20 AS runtime
WORKDIR /app
COPY --from=build /out/api ./api
EXPOSE 8080
ENTRYPOINT ["/app/api"]

FROM golang:1.22-alpine AS dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/air-verse/air@v1.52.3
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
