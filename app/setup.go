package app

import (
	"revel-project/app/models"
	"gorm.io/driver/postgres"
	"github.com/revel/revel"
	"gorm.io/gorm"
	"fmt"
)

// -------- DB Setup --------

type DBConn struct {
	host				string
	user				string
	password		string
	name				string
	port				int

	db					*gorm.DB
}


func(this *DBConn) FireUp() (error) {
	revel.AppLog.Info("Firing Up DB!")

	connstring := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		this.host, this.user, this.password, this.name, this.port,
	)

	db, err := gorm.Open(postgres.Open(connstring), &gorm.Config{})

	if err != nil {
		return err
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	this.db = db
	this.migrate()

	revel.AppLog.Info("DB all fired up!")
	return nil
}


func(this DBConn) CoolDown() {
	revel.AppLog.Info("Cooling down DB")

	sqlDB, sqlErr := this.db.DB()
	
	if sqlErr != nil {
		panic(sqlErr)
	}

	closeErr := sqlDB.Close()

	if closeErr != nil {
		panic(closeErr)
	}
}


func(this *DBConn) migrate() {

	// revel.AppLog.Warn("Dropping Users table")
	// db.Migrator().DropTable(&models.User{})

	revel.AppLog.Info("Migrating...")
	// add DB table models here
	this.db.AutoMigrate(
		&models.User{},
	)
}

// ---- setters ----

func(this *DBConn) SetHost(host string) {
	this.host = host
}

func(this *DBConn) SetUser(user string) {
	this.user = user
}

func(this *DBConn) SetPassword(password string) {
	this.password = password
}

func(this *DBConn) SetName(name string) {
	this.name = name
}

func(this *DBConn) SetPort(port int) {
	this.port = port
}

// ---- Getters ----

func(this DBConn) GetDB() *gorm.DB {
	return this.db
}
