build:
	@go build -o bin/echoes ./main.go
	@cp .env bin

run: build
	@./bin/echoes

test:
	@go test

lint:
	@golangci-lint run

docker:
	@docker build -t yosaa5782/echoes .
	@docker run --rm --name echoes-server -p 5000:5000 --network echoes-net -d yosaa5782/echoes

postgres:
	@docker run --rm \
		--name postgres-echoes \
		-p 5432:5432 \
		-e POSTGRES_USER=user \
		-e POSTGRES_PASSWORD=1234 \
		-v /postgres-volume:/var/lib/postgresql/data \
		--network echoes-net \
		-d postgres
	@docker run --rm \
		--name adminer-echoes \
		-p 5050:8080 \
		--network echoes-net \
		-d adminer

create-db:
	@docker exec -it postgres-echoes createdb --username=user --owner=user echoesdb

create-network:
	@docker network create -d bridge echoes-net

mig:
	@migrate create -ext sql -dir migrations -seq migration_name

migrate-up:
	@migrate -path migrations -database "postgres://user:1234@localhost:5432/echoesdb?sslmode=disable" -verbose up

migrate-down:
	@migrate -path migrations -database "postgres://user:1234@localhost:5432/echoesdb?sslmode=disable" -verbose down

migrate-fix:
	@migrate -path migrations -database "postgres://user:1234@localhost:5432/echoesdb?sslmode=disable" force VERSION

redis:
	@docker run --rm \
		--name redis-echoes \
		-p 6379:6379 \
		-d redis

redis-cli:
	@docker exec -it redis-echoes redis-cli