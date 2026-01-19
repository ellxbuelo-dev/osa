# syntax=docker/dockerfile:1

# ====== BUILD STAGE ======
FROM golang:1.25-alpine AS build

WORKDIR /src

# Сначала зависимости (кэш)
COPY go.mod go.sum ./
RUN go mod download

# Потом код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api .

# ====== RUNTIME STAGE ======
FROM gcr.io/distroless/static:nonroot
WORKDIR /

COPY --from=build /out/api /api

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/api"]
