run:
	go run cmd/api/main.go

migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations \
	-database "postgres://admin:admin@localhost:5432/security_monitor?sslmode=disable" \
	up

migrate-down:
	migrate -path migrations \
	-database "postgres://admin:admin@localhost:5432/security_monitor?sslmode=disable" \
	down