# DB Migrate

## How to use

enter your DB container

```bash
docker exec -it container_name psql -U xxx
```

create new migration file

```bash
# create a new table
./migrate create -ext sql -dir db/migrations create_xxx_table
# run db migrate
./migrate -database "postgres://db_url" -path db/migrations up
# run db rollback
./migrate -database "postgres://db_url" -path db/migrations down
```

## SQL Note

- How does Postgres handle text column?
  - Text doesn't have actual limit. The actual limit is constrained by the storage size of the table.
  - Postgres can handle large data values efficiently by compressing and store them outside the table.
- Error: Dirty database version 20231219090208. Fix and force version.
  - Why?
    - Before a migration runs, each db sets a dirty flag. Execuation stops if a migration fails and dirty state persists.
  - How to solve it?
    1. manually fix the migration file (sql syntax errors...)
    2. force the expected version -> (`go run main.go migrate -force=20231219090208`)
