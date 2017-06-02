package models

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	var err error
	driver_name := "mysql"
	if driver_name == "" {
		log.Fatal("Invalid driver name")
	}
	dsn := "root:@/simple_development?charset=utf8&parseTime=True&loc=Local"
	if dsn == "" {
		log.Fatal("Invalid DSN")
	}
	db, err = sqlx.Connect(driver_name, dsn)
	if err != nil {
		log.Fatal(err)
	}
}
