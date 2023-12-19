package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const (
  migrationPaths = "file://db/migrations"
)

var (
  createDBCmd = flag.NewFlagSet("createdb", flag.ExitOnError)
  createDBName = createDBCmd.String("name", "", "given the name of the database")

  migrateCmd = flag.NewFlagSet("migrate", flag.ExitOnError)
  forceVersion = migrateCmd.Int("force", 0, "force a specific version")

  rollbackCmd = flag.NewFlagSet("rollback", flag.ExitOnError)
  rollbackStep = rollbackCmd.Int("step", 1, "rollback 1 step by default")
)

func main() {
  flag.Parse()
  
  if len(flag.Args()) < 1 {
    log.Fatal("only accept commands: createdb, migrate, or rollback") 
  }

  if err := godotenv.Load(); err != nil {
		log.Fatal(err) 
	}
  
  switch flag.Args()[0] {
  case "createdb":
    createDBCmd.Parse(flag.Args()[1:])
    if err := runCreateDB(*createDBName); err != nil {
      log.Fatal(err)
    }
  case "migrate":
    migrateCmd.Parse(flag.Args()[1:])
     if err := runMigrate(*forceVersion); err != nil {
      log.Fatal(err)
    }
  case "rollback":
    rollbackCmd.Parse(flag.Args()[1:])
    if err := runRollback(*rollbackStep); err != nil {
      log.Fatal(err)
    }
  default:
    log.Fatal("unknown command. Please use createdb, migrate, or rollback")
  }
}

func runCreateDB(name string) error {
  if name == "" {
    return errors.New("no database name")
  }
  db, err := sql.Open("postgres", postgresUrl())
  if err != nil {
    return err
  }
  defer db.Close()
  createDbSql := fmt.Sprintf("CREATE DATABASE %s", name)
  if _, err := db.Exec(createDbSql); err != nil {
    return err
  } 
  return nil
}

func runMigrate(version int) error {
  m, err := migrate.New(
    migrationPaths,
    postgresDbUrl(),
   )
  if err != nil {
    return err
  }
  if version > 0 {
    if err := m.Force(version); err != nil {
      return err
    }
  } else {
    if err := m.Up(); err != nil {
      return err
    }
  }
  m.Close()
  return nil
}

func runRollback(step int) error {
  m, err := migrate.New(
    migrationPaths,
    postgresDbUrl(),
   )
  if err != nil {
    return err
  }
  if err := m.Steps(-step); err != nil {
    return err
  }
  m.Close()
  return nil
}

func postgresUrl() string {
  return fmt.Sprintf(
    "postgres://%s:%s@%s:%v?sslmode=disable",
    os.Getenv("username"),
    os.Getenv("password"),
    os.Getenv("host"),
    os.Getenv("port"),
  )
}

func postgresDbUrl() string {
  return fmt.Sprintf(
    "postgres://%s:%s@%s:%s/%s?sslmode=disable",
    os.Getenv("username"),
    os.Getenv("password"),
    os.Getenv("host"),
    os.Getenv("port"),
    os.Getenv("dbname"),
  )
}
