package main

import (
	"database/sql"
	"log"

	"github.com/gitaepark/carrot-market/controller"
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/loader"
	"github.com/gitaepark/carrot-market/service"
	"github.com/gitaepark/carrot-market/token"
	"github.com/gitaepark/carrot-market/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannto connect to db: ", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	runGinServer(config, store)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance: ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	tokenMaker, err := token.NewJWTMaker(config.JWTSecret)
	if err != nil {
		log.Fatal("cannot create token maker: ", err)
	}

	service := service.NewService(config, tokenMaker, store)
	controller, err := controller.NewController(service)
	if err != nil {
		log.Fatal("cannot create controller: ", err)
	}

	server, err := loader.NewServer(controller)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
