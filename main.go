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

  rollbackCmd = flag.NewFlagSet("rollback", flag.ExitOnError)
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
     if err := runMigrate(); err != nil {
      log.Fatal(err)
    }
  case "rollback":
    rollbackCmd.Parse(flag.Args()[1:])
    
    if err := runRollback(); err != nil {
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

func runMigrate() error {
  m, err := migrate.New(
    migrationPaths,
    postgresDbUrl(),
   )
  if err != nil {
    return err
  }
  if err := m.Up(); err != nil {
    return err
  }
  m.Close()
  return nil
}

func runRollback() error {
  m, err := migrate.New(
    migrationPaths,
    postgresDbUrl(),
   )
  if err != nil {
    return err
  }
  if err := m.Down(); err != nil {
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
