.PHONY: help build up down restart logs k6 migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build        - Build all Docker images"
	@echo "  make up           - Start all services (PostgreSQL, Redis, Backend, Frontend)"
	@echo "  make down         - Stop all services"
	@echo "  make restart      - Restart all services"
	@echo "  make logs         - View logs for all running services"
	@echo "  make k6           - Execute k6 load test (1,000 VUs benchmark)"

build:
	docker compose build

up:
	docker compose up --build -d

down:
	docker compose down

restart:
	docker compose restart

logs:
	docker compose logs -f

k6:
	docker run --rm -i --network="host" grafana/k6 run - < k6/load_test.js
