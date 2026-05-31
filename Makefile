.PHONY: up down test lint

up:
	docker compose up --build

down:
	docker compose down

test:
	cd runfactory-core && go test ./... && go build ./...
	cd runfactory-dashboard && npm run build

lint:
	cd runfactory-dashboard && npm run lint
