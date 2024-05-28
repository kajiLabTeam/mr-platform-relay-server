-include .env

up:
	docker compose build && docker compose up -d

logs:
	docker compose logs -f

down:
	docker compose down

go:
	docker exec -it $(GO_CONTAINER_HOST) /bin/sh