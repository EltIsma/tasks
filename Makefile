migrate-create:
	migrate create -ext sql -dir migrations -seq init

migrate-up:
	migrate -path migrations -database $(url) up

migrate-down:
	migrate -path migrations -database $(url) down

compose:
	docker-compose up --remove-orphans --build

.PHONY: swagger
swagger:
	swag init -g internal/ports/httpServer/router.go -o docs/swagger --parseInternal

create-kafka-topic:
	kafka-scripts/create-topic.sh

delete-kafka-topic:
	kafka-scrits/delete-topic.sh