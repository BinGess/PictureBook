FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://proxy.golang.org,direct
RUN go mod tidy
RUN go build -o server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/server /app/server
COPY --from=build /app/web /app/web
ENV CORS_ALLOWED_ORIGINS=*
RUN mkdir -p /app/uploads
EXPOSE 8080
CMD ["/app/server"]