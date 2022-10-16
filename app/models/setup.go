package models

import (
	"fmt"
	"gorm.io/gorm"
)

// TEMPORARY, definitely better ways to structure DB setup... :)

func FireUp(user string, password string, dbname string, port string) (*gorm.DB, error) {
	connstring := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable",
		user, password, dbname, port
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return db, err
	}

	migrate(db)

	return db, nil
}

func CoolDown(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}

func migrate(db *gorm.DB) {
	// add DB table models here
	db.AutoMigrate(
		&User{},
	)

	// removing column/tables/etc goes down here, since AutoMigrate doesn't remove data...?

}