.PHONY: run

run:
	sudo docker compose --file deploy/docker-compose.yml down
	sudo docker compose -f deploy/docker-compose.yml build --no-cache
	sudo docker compose --file deploy/docker-compose.yml up --scale rate-limiter=2
