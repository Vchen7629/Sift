core_up:
	docker compose --profile core up -d

core_down:
	docker compose --profile core down -v

observability_up:
	docker compose --profile observability up -d

observability_down:
	docker compose --profile observability down -v

observability_restart: 
	docker compose --profile observability down -v && docker compose --profile observability up -d