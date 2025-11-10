FROM golang:1.25.3-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates && update-ca-certificates

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/. .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o /out/migrator ./cmd/migrator/migrator.go

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/migrator .
COPY backend/service/db/migrations ./migrations
USER nonroot:nonroot
ENTRYPOINT ["./migrator"]