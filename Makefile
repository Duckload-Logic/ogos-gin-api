include .env
export

# Define db url
DB_URL=mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

# Desc: To create a new migration
# Usage: make migration name=create_users
migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

# Desc: To create insert seed data migration
# Usage: make seed name=insert_initial_data
seed:
	migrate create -ext sql -dir ./seeds -seq $(name)

# Desc: To create fake data seed migration
# Usage: make fake
fakes:
	go run cmd/faker/faker.go

# Desc: Fetch PSGC data from API and save to JSON (run once, commit the file)
# Usage: make psgc-fetch
psgc-fetch:
	go run cmd/psgc/psgc.go

# Desc: Seed locations from psgc_data.json into the database
# Usage: make locations
locations:
	go run cmd/locations/locations.go

# Desc: To refresh database with cli
# Usage: make migrate-up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

# Desc: To undo database migrations
# Usage: make migrate-down
migrate-down:
	migrate -path migrations -database "$(DB_URL)" drop -f

# Desc: To apply seed data to database
# Usage: make seed-up
seed-up:
	migrate -path seeds -database \
	"$(DB_URL)?x-migrations-table=seed_migrations" up

refresh: migrate-down migrate-up seed-up locations

# Desc: To generate swagger docs
# Usage: make swagger-internal
swagger-internal:
	swag init -g main.go \
	--parseDependency --parseInternal \
	--dir ./cmd/api,\
	./internal/features/auth,\
	./internal/features/users,\
	./internal/features/appointments,\
	./internal/features/excuseslips,\
	./internal/features/students
	--output ./docs/internal \
    --instanceName internal

swagger-external:
	swag init -g main.go \
	--parseDependency --parseInternal \
	--dir ./cmd/api,\
	./internal/features/students/external \
	--output ./docs/external \
	--instanceName external