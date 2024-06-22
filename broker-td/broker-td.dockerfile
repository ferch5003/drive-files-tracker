# Build tiny docker image.
FROM alpine:latest

RUN mkdir /app

COPY go.mod /
COPY .env /
COPY brokerTDApp /app

CMD ["/app/brokerTDApp"]