package main

import (
	"database/sql"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mobyos/admin-app/server/db"
	"github.com/mobyos/admin-app/server/types"
)

var descriptor []byte = []byte(`

name: app1
description: "This app is amazing and it do magical stuff"
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

	b, err := sql.Open("sqlite3", "../db/mobyos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	sqlStmt := `create table application (id text not null primary key, name text, descriptor blob);`
	_, err = b.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	desc := types.AppDescriptor{}
	err = yaml.Unmarshal(descriptor, &desc)
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateApplication(desc)

	if err != nil {
		log.Fatal(err)
	}

}
