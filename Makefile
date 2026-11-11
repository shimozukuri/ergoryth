include .env
export

proto-lint:
	@buf lint

proto:
	@make proto-lint && \
	rm -rf .\backend\internal\generated && \
	buf generate

backend-run:
	@cd backend && \
	go mod tidy && \
	go run cmd\ergoryth\main.go