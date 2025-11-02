FROM golang:1.25.3-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates && update-ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/server .
USER nonroot:nonroot
ENTRYPOINT ["./server"]