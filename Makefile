up:
	docker compose up --build

dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml --profile dev up --build

down:
	docker compose down -v

test:
	docker compose run --rm api-dev go test ./...

tidy:
	docker compose run --rm api-dev go mod tidy
