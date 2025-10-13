FROM golang:1.25.2-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates && update-ca-certificates
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -trimpath -o /out/app .

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates
WORKDIR /app
COPY --from=build /out/app /app/app
USER appuser
ENTRYPOINT ["/app/app"]
