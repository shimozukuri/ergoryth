include .env
export

proto-lint:
	@buf lint

proto:
	@make proto-lint
	buf generate

backend-run:
	@cd backend && \
	go mod tidy && \
	go run cmd/ergoryth/main.go