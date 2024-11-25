FROM golang:latest AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o telegramApp ./cmd/bot

FROM alpine:latest AS build-release-stage

# Installs latest Chromium package.
RUN apk update && apk upgrade && apk add --no-cache bash git && apk add --no-cache chromium
RUN echo @edge http://nl.alpinelinux.org/alpine/edge/community >> /etc/apk/repositories \
    && echo @edge http://nl.alpinelinux.org/alpine/edge/main >> /etc/apk/repositories \
    && apk add --no-cache \
    harfbuzz@edge \
    nss@edge \
    freetype@edge \
    ttf-freefont@edge \
    && rm -rf /var/cache/* \
    && mkdir /var/cache/apk
CMD chromium-browser --headless --disable-gpu --safebrowsing-disable-auto-update --disable-sync --disable-default-apps \
    --hide-scrollbars --mute-audio --no-first-run --no-sandbox

WORKDIR /app
COPY --from=build-stage /app/telegramApp /app/telegramApp
ENTRYPOINT ["/app/telegramApp"]