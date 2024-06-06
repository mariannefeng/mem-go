set dotenv-load

# brings up postgres db in a docker container
db:
    docker-compose up -d

# runs backend with live reload on 8080
run: db
    air

# runs database migration
db-migrate: db
    GOOSE_DRIVER=$DB_DRIVER GOOSE_DBSTRING=$DB_STRING GOOSE_MIGRATION_DIR=$MIGRATION_DIR goose up

# returns database migration status
db-status: db
    GOOSE_DRIVER=$DB_DRIVER GOOSE_DBSTRING=$DB_STRING goose status

db-generate:
    sqlc generate