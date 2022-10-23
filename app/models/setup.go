package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/revel/revel"
)

// TEMPORARY, definitely better ways to structure DB setup... :)

func FireUp(user string, password string, dbname string, port string) (*gorm.DB, error) {
	connstring := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable",
		user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(connstring), &gorm.Config{})

	if err != nil {
		return db, err
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	migrate(db)

	return db, nil
}

func CoolDown(db *gorm.DB) {
	sqlDB, sqlErr := db.DB()
	
	if sqlErr != nil {
		panic(sqlErr)
	}

	closeErr := sqlDB.Close()

	if closeErr != nil {
		panic(closeErr)
	}
}

func migrate(db *gorm.DB) {

	// revel.AppLog.Warn("Dropping Users table")
	// db.Migrator().DropTable(&User{})

	revel.AppLog.Info("Migrating...")
	// add DB table models here
	db.AutoMigrate(
		&User{},
	)
}