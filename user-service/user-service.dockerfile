# Build tiny docker image.
FROM alpine:latest

RUN mkdir /app

COPY go.mod /
COPY .env /
COPY userServiceApp /app

CMD ["/app/userServiceApp"]