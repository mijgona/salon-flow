.PHONY: build run test vet lint clean

build:
	go build -o bin/salon-crm ./cmd/app

run:
	go run ./cmd/app

test:
	go test ./... -v

test-domain:
	go test ./internal/core/domain/... -v

vet:
	go vet ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

migrate-up:
	@echo "Run migrations manually: psql -f migrations/001_clients.sql && ..."

docker-build:
	docker build -t salon-crm .

docker-run:
	docker run -p 8080:8080 salon-crm
