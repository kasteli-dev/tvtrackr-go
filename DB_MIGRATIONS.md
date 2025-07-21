# Database Migrations with golang-migrate

## How to run migrations

1. Make sure you have the `migrate` CLI installed with Postgres support:

   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

2. Run migrations (replace the URL if needed):

   ```bash
   migrate -path ./migrations -database "postgres://tvtrackr_postgresql_user:ZOQvXIBVqBveNKJFrSPadIO4nMaAjbH0@dpg-d1v51mbipnbc73apcbjg-a.frankfurt-postgres.render.com/tvtrackr_postgresql" up
   ```

Or, if you have `POSTGRES_URL` in your `.env`:

   ```bash
   migrate -path ./migrations -database "$POSTGRES_URL" up
   ```

## Creating new migrations

To create a new migration file:

```bash
migrate create -ext sql -dir ./migrations -seq name_of_migration
```

## Useful commands

- Run all up migrations: `migrate -path ./migrations -database "$POSTGRES_URL" up`
- Rollback last migration: `migrate -path ./migrations -database "$POSTGRES_URL" down 1`
- Check current version: `migrate -path ./migrations -database "$POSTGRES_URL" version`

## Docs
- https://github.com/golang-migrate/migrate
