package main

import (
	"database/sql"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mobyos/mobyos-admin-app/server/db"
	"github.com/mobyos/mobyos-admin-app/server/types"
)

var descriptors [][]byte = [][]byte{[]byte(`

name: "VLC"
description: "VLC is a free and open source cross-platform multimedia player and framework that plays most multimedia files as well as DVDs, Audio CDs, VCDs, and various streaming protocols."
icon_url: "http://i.utdstc.com/icons/256/vlc-media-player-1-0-5.png"
remote_path: "/mobile.html"
services:
  app:
    image: jess/vlc
    ui: true
    sound: true
    command: [--control=http, --http-host=0.0.0.0, --http-port=9000, --http-password=lalala]
    ports:
        - "9000"

`), []byte(`
name: "Kodi"
description: "Kodi, the one and only media center"
icon_url: "http://www.homemediatech.net/wp-content/uploads/2015/11/kodi-logo.png"
services:
  app:
    image: marcosnils/kodi
    ui: true
    sound: true

`)}

func main() {
	os.Remove("../db/mobyos.db")

	b, err := sql.Open("sqlite3", "../db/mobyos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	sqlStmt := `create table application (id text not null primary key, descriptor blob);`
	_, err = b.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	for _, descriptor := range descriptors {
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

}
