FROM golang:latest AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o userServiceApp ./cmd/api

FROM alpine:latest AS build-release-stage
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build-stage /app/userServiceApp /app/userServiceApp
ENTRYPOINT ["/app/userServiceApp"]