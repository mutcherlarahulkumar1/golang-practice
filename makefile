
DB_URL=postgresql://neondb_owner:npg_tGkoOZ5yB9nI@ep-summer-smoke-a4sqvd3z-pooler.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require
MIGRATIONS_PATH=sqlconnect/migrations

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# Rollback last migration
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

# Rollback all migrations
migrate-reset:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

# Show current migration version
migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version
