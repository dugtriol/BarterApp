include .env
GOLINT := $(GOPATH)/bin/golint

.PHONY: migrate-down
migrate-down:
	migrate -database '$(POSTGRES_CONN_LOCAL)?sslmode=disable' -path migrations down


.PHONY: migrate-up
migrate-up:
	migrate -database '$(POSTGRES_CONN_LOCAL)?sslmode=disable' -path migrations up

.PHONY: compose-up
compose-up:
	docker-compose up

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: install-lint
install-lint:
	go install golang.org/x/lint/golint@latest

.PHONY: lint
lint:
	$(GOLINT) ./...

.PHONY: generate
generate:
	go run github.com/99designs/gqlgen generate