ENV ?= .env
include $(ENV)
export

# Define db url
DB_URL=mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)
ifeq ($(DB_TLS),true)
DB_URL:=$(DB_URL)?tls=true
endif

# Desc: To create a new migration
# Usage: make migration name=create_users
migration:
	migrate create -ext sql -dir ./scripts/migrations -seq $(name)

# Desc: To create insert seed data migration
# Usage: make seed name=insert_initial_data
seed:
	migrate create -ext sql -dir ./scripts/seeds -seq $(name)

# Desc: To create fake data seed migration
# Usage: make fake
fakes:
	go run ./scripts/faker

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
	migrate -path scripts/migrations -database "$(DB_URL)" up

# Desc: To undo database migrations
# Usage: make migrate-down
migrate-down:
	migrate -path scripts/migrations -database "$(DB_URL)" down

# Desc: To undo all database migrations
# Usage: make migrate-reset
migrate-reset:
	migrate -path scripts/migrations -database "$(DB_URL)" drop -f

# Desc: To apply seed data to database
# Usage: make seed-up
seed-up:
	migrate -path scripts/seeds -database \
	"$(DB_URL)$(if $(filter true,$(DB_TLS)),&,?)x-migrations-table=seed_migrations" up

refresh: migrate-reset migrate-up seed-up locations

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

compose-up:
	docker-compose --env-file $(ENV) up --build

compose-prod:
	DOCKER_STAGE=prod docker-compose --env-file $(ENV) up --build -d
