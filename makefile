.PHONY: migrate-create migrate-up migrate-down help

migrate-create:
	@echo "Creating migrations..."
	@if [ -z "$(name)" ]; then \
		echo "Please provide a migration name." && exit 1; \
	fi
	migrate create -ext sql -dir cmd/migrate/migrations $(name)
	@echo "Migrations created successfully!"

migrate-up:
	@echo "Applying migrations..."
	go run cmd/migrate/main.go up
	@echo "Migrations applied successfully!"

migrate-down:
	@echo "Reverting migrations..."
	go run cmd/migrate/main.go down
	@echo "Migrations reverted successfully!"

help:
	@echo "Invalid option. Use:"
	@echo "  make migrate-create name=<migration_name> - to create a migration"
	@echo "  make migrate-up - to apply migrations"
	@echo "  make migrate-down - to revert migrations"
