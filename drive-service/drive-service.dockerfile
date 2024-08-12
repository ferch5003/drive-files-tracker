FROM golang:latest as build-stage
WORKDIR /app
COPY ./config/client_secret.json /config
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o driveServiceApp ./cmd/api

FROM alpine:latest as build-release-stage
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build-stage /app/driveServiceApp /app/driveServiceApp
ENTRYPOINT ["/app/driveServiceApp"]