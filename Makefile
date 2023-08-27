include .env

migrate-up:
	migrate -path migrations -database "postgresql://$(PGUSER):$(PGPASSWORD)@$(PGHOST):$(PGPORT)/$(PGDATABASE)?sslmode=$(PGSSLMODE)" -verbose up

migrate-down:
	migrate -path migrations -database "postgresql://$(PGUSER):$(PGPASSWORD)@$(PGHOST):$(PGPORT)/$(PGDATABASE)?sslmode=$(PGSSLMODE)" -verbose down

createdb:
	docker run --name avitodb -e POSTGRES_DB=$(PGDATABASE) -e POSTGRES_USER=$(PGUSER) -e POSTGRES_PASSWORD=$(PGPASSWORD) -p 5432:5432 -d postgres

run-docker:
	docker-compose up 

run-up:
	go run cmd/main.go

test:
	go test  ./...

.PHONY: migrate-down migrate-up run-up
