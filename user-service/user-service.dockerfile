# Build tiny docker image.
FROM alpine:latest

RUN apk add --no-cache tzdata

RUN mkdir /app

COPY go.mod /
COPY .env /
COPY userServiceApp /app

CMD ["/app/userServiceApp"]