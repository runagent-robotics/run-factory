.PHONY: up down test lint

up:
	docker compose up --build

down:
	docker compose down

test:
	cd runfactory-core && go test ./... && go build ./...
	cd runfactory-dashboard && npm run build

lint:
	cd runfactory-core && go vet ./...
	cd runfactory-dashboard && npm run lint
