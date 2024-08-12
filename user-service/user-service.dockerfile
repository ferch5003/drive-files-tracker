FROM golang:latest as build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o userServiceApp ./cmd/api

FROM alpine:latest as build-release-stage
WORKDIR /app
COPY --from=build-stage /app/userServiceApp /app/userServiceApp
ENTRYPOINT ["/app/brokerTDApp"]