service = lazer
container = lazer-dev

bash:
	docker exec -it $(container) bash

build-dev:
	docker-compose build $(service)

up:
	docker-compose up -d

log:
	docker logs -f --tail=100 $(container)

down:
	docker-compose down

fresh: down build-dev up
