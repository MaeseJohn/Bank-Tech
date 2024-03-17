package migrations

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Migrations *migrate.Migrate

func CreateDataBase() {
	// Creating db with migrations
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	var err error
	pg := postgres.Open("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	gdb, _ := gorm.Open(pg, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	db, err := gdb.DB()
	if err != nil {
		panic(err)
	}
	config := migratePostgres.Config{}
	driver, err := migratePostgres.WithInstance(db, &config)
	if err != nil {
		panic(err)
	}

	Migrations, err = migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"postgres",
		driver)
	if err != nil {
		log.Fatal(err)
	}
	err = Migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
