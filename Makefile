include .env
export

# Define db url
DB_URL=mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

# Desc: To create a new migration
# Usage: make migration name=create_users
migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

# Desc: To refresh database with cli
# Usage: make migrate-up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

# Desc: To undo database migrations
# Usage: make migrate-down
migrate-down:
	migrate -path migrations -database "$(DB_URL)" drop -f