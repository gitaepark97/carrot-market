package main

import (
	"log"

	"github.com/gitaepark/carrot-market/controller"
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/loader"
	"github.com/gitaepark/carrot-market/service"
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

	conn, err := loader.ConnectDB(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannto connect to db: ", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	service, err := service.NewService(config, store)
	if err != nil {
		log.Fatal("cannot create service: ", err)
	}

	runGinServer(service, config.HTTPServerAddress)
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

func runGinServer(service *service.Service, address string) {
	controller, err := controller.NewController(service)
	if err != nil {
		log.Fatal("cannot create controller: ", err)
	}

	server, err := loader.NewServer(controller)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
