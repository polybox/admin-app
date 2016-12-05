package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Remove("../db/mobyos.db")

	b, err := sql.Open("sqlite3", "../db/mobyos.db")
	defer b.Close()
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `create table application (id text not null primary key, descriptor blob);`
	_, err = b.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

}
