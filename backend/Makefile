run:
	docker compose build
	docker compose up -d

stop:
	docker compose down
	
test:
	@echo "Building Project..."
	@docker compose -f docker-compose.test.yml build -q
	@echo "Starting Tests..."
	@if docker compose -f docker-compose.test.yml up --abort-on-container-exit --no-attach db; then \
		echo "Tests Passed!"; \
		docker compose down; \
		exit 0; \
	else \
		echo "Tests Failed!"; \
		docker compose down; \
		exit 1; \
	fi

wipe:
	docker volume rm backend_postgresql
	docker volume rm backend_uploads