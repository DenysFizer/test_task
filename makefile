migrations-up:
	migrate -database 'postgres://admin:admin@localhost:5432/db?sslmode=disable' -path migrations up
migrations-down:
	migrate -database 'postgres://admin:admin@localhost:5432/db?sslmode=disable' -path migrations down