all_up: core_up observability_up
all_down: core_down core_up
all_restart: core_restart observability_restart

core_up:
	docker compose --profile core up -d
core_down:
	docker compose --profile core down -v
core_restart: core_down core_up


observability_up:
	docker compose --profile observability up -d
observability_down:
	docker compose --profile observability down -v
observability_restart: observability_down observability_up