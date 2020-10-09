package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Db *sqlx.DB
}

// Connect ...
func Connect(host, port, uname, pass, dbname string) (*DB, error) {
	var db = &DB{}
	dbSource := fmt.Sprintf(
		"root:%s@tcp(%s:%s)/%s?charset=utf8",
		pass,
		host,
		port,
		dbname,
	)
	d, err := sqlx.Connect("mysql", dbSource)
	db.Db = d

	return db, err
}
