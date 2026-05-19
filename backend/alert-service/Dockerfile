FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/service ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/service /app/service
EXPOSE 8084
CMD ["/app/service"]
