run:
	docker compose build
	docker compose up -d

stop:
	docker compose down

unittest:
	go test ./... -tags=unittest -cover

integrationtest:
	docker compose -f docker-compose.test.yml build
	-docker compose -f docker-compose.test.yml up --abort-on-container-exit
	docker compose down