build: build_broker build_auth build_logger
	docker compose down
	docker compose up --build -d

down:
	docker compose down

build_broker:
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o brokerService ./cmd/api

build_auth:
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o authService ./cmd/api

build_logger:
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o loggerService ./cmd/api

build_image: build_broker build_logger build_auth build_broker_service build_logger_service build_auth_service

build_logger_service:
	cd ../logger-service && docker build -f logger-service.dockerfile -t muzaffertolgayakar/logger-service:2.0.0 . && docker push muzaffertolgayakar/logger-service:2.0.0
build_broker_service:
	cd ../broker-service && docker build -f broker-service.dockerfile -t muzaffertolgayakar/broker-service:2.0.0 . && docker push muzaffertolgayakar/broker-service:2.0.0
build_auth_service:
	cd ../auth-service && docker build -f auth-service.dockerfile -t muzaffertolgayakar/auth-service:2.0.0 . && docker push muzaffertolgayakar/auth-service:2.0.0