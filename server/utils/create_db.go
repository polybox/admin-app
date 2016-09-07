package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/twinj/uuid"
)

var descriptor []byte = []byte(`

name: app1
services:
  app:
    image: nginx
  web:
    image: nginx
    ports:
      - "8000:8000"

`)

func main() {
	os.Remove("../db/mobyos.db")

	db, err := sql.Open("sqlite3", "../db/mobyos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `create table installation (id text not null primary key, application_id text, descriptor blob);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into installation(id, application_id, descriptor) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	app1 := uuid.NewV4()
	_, err = stmt.Exec(app1.String(), "app1", descriptor)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100; i++ {
	}
	tx.Commit()
}
