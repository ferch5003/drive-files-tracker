# Build tiny docker image.
FROM alpine:latest

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

RUN mkdir /app

COPY go.mod /
COPY .env /
COPY telegramApp /app

CMD ["/app/telegramApp"]