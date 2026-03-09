FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -tags timetzdata -o bin/api ./cmd/api

FROM alpine:3.21
RUN apk --no-cache add ca-certificates && \
    addgroup -S usergroup && adduser -S user -G usergroup
WORKDIR /app
COPY --from=builder /app/bin/api .
RUN chown user:usergroup /app/api
USER user
EXPOSE 8080
CMD [ "./api" ]
