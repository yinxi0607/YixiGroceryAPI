.PHONY: tidy run-user

run-user:
	go run user-service/main.go

run-api:
	go run api-gateway/main.go