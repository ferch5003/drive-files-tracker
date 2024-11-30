FROM golang:alpine AS build-stage
WORKDIR /app

RUN apk update
RUN apk add  tesseract-ocr-dev \
             g++ \
             musl-dev \
             leptonica-dev \
             tesseract-ocr-data-eng \
             tesseract-ocr-data-spa
ENV GOSSERACT_CPPSTDERR_NOT_CAPTURED=1

COPY ./config/client_secret.json /config
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=1 go build -o driveServiceApp ./cmd/api
RUN ["chmod", "+x", "/app/driveServiceApp"]

FROM alpine:latest AS build-release-stage

RUN apk add --no-cache tzdata

RUN apk update
RUN apk add  tesseract-ocr-dev \
             g++ \
             musl-dev \
             leptonica-dev \
             tesseract-ocr-data-eng \
             tesseract-ocr-data-spa
ENV GOSSERACT_CPPSTDERR_NOT_CAPTURED=1

WORKDIR /app
COPY --from=build-stage /app/driveServiceApp /app/driveServiceApp
COPY --from=build-stage /app/config/client_secret.json /app/config/client_secret.json
COPY --from=build-stage /app/go.mod /app/go.sum ./
RUN mkdir -m 777 tmp
ENTRYPOINT ["/app/driveServiceApp"]