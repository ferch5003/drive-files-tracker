TELEGRAM_BINARY=telegramApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_telegram
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_telegram:
	@echo "Building Telegram Bots binary..."
	cd ../telegram-bot-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TELEGRAM_BINARY} ./cmd/bot
	@echo "Done!"

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_BINARY}"
	@echo "Stopped front end!"