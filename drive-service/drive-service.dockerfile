FROM golang:latest AS build-stage
WORKDIR /app
COPY ./config/client_secret.json /config
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o driveServiceApp ./cmd/api
RUN ["chmod", "+x", "/app/driveServiceApp"]

FROM alpine:latest AS build-release-stage
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build-stage /app/driveServiceApp /app/driveServiceApp
COPY --from=build-stage /app/config/client_secret.json /app/config/client_secret.json
COPY --from=build-stage /app/go.mod /app/go.sum ./
USER root
RUN mkdir -m 777 tmp
ENTRYPOINT ["/app/driveServiceApp"]