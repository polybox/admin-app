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
remote_path: "/"
services:
  app:
    image: marcosnils/kodi
    ui: true
    sound: true
    ports:
        - "8080"

`), []byte(`
name: "Motion"
description: "See everything"
icon_url: "https://cdn4.iconfinder.com/data/icons/technology-devices-1/500/security-camera-128.png"
remote_path: "/"
services:
  app:
    image: surround/rpi-motion-mmal
    sound: true
    ports:
        - "8081"

`), []byte(`
name: "Retropie"
description: "Play your favourite Arcade, home-console, and classic PC games with the minimum set-up."
icon_url: "https://retroresolution.files.wordpress.com/2016/03/retropie_logo_300x300.png"
remote_path: "/"
services:
  app:
    image: retropie
    sound: true
    input: true
    ports:
        - "8080"

`), []byte(`
name: "Hotspot"
description: "Share your internet connection with your guests"
icon_url: "http://www.montclair-hostel.com/wp-content/uploads/2015/03/wifi.png"
remote_path: "/"

`), []byte(`
name: "Secure nav"
description: "Navigate through tor router securely"
icon_url: "https://upload.wikimedia.org/wikipedia/commons/7/73/Tor_logo-1.png"
remote_path: "/"
`), []byte(`
name: "Media server"
description: ""
icon_url: "https://upload.wikimedia.org/wikipedia/commons/7/73/Tor_logo-1.png"
remote_path: "/"
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
