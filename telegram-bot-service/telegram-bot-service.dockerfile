# Build tiny docker image.
FROM alpine:latest

RUN mkdir /app

COPY go.mod /
COPY .env /
COPY telegramApp /app

CMD ["/app/telegramApp"]