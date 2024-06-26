BROKER_BINARY=brokerTDApp
TELEGRAM_BINARY=telegramApp
USER_BINARY=userServiceApp
DRIVE_BINARY=driveServiceApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_user build_telegram build_drive
	@echo "Stopping docker images (if running...)"
	docker-compose down --remove-orphans
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-td && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_telegram:
	@echo "Building Telegram Bots binary..."
	cd ../telegram-bot-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TELEGRAM_BINARY} ./cmd/bot
	@echo "Done!"

## build_broker: builds the user binary as a linux executable
build_user:
	@echo "Building User Service binary..."
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd/api
	@echo "Done!"

## build_broker: builds the drive binary as a linux executable
build_drive:
	@echo "Building Drive Service binary..."
	cd ../drive-service && env GOOS=linux CGO_ENABLED=0 go build -o ${DRIVE_BINARY} ./cmd/api
	@echo "Done!"