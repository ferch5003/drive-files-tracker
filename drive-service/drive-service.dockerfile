# Build tiny docker image.
FROM alpine:latest

RUN apk add --no-cache tzdata

RUN mkdir /config

COPY ./config/client_secret.json /config

RUN mkdir /app

COPY go.mod /
COPY .env /
COPY driveServiceApp /app

CMD ["/app/driveServiceApp"]