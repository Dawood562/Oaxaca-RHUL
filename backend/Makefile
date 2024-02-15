run:
	docker compose build
	docker compose up -d

stop:
	docker compose down

test:
	docker compose -f docker-compose.test.yml build
	-docker compose -f docker-compose.test.yml up --abort-on-container-exit --no-attach db
	docker compose down

wipe:
	docker volume rm backend_postgresql