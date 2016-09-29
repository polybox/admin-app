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

name: "Spotify"
description: "Music for everyone"
icon_url: "https://play.spotify.edgekey.net/site/0298183/images/favicon.png"
remote_path: "/musicbox_webclient"
services:
  app:
    image: mopidy
    sound: true
    ports:
        - "6680"

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
name: "Hotspot"
description: "Share your internet connection with your guests"
icon_url: "http://www.montclair-hostel.com/wp-content/uploads/2015/03/wifi.png"
remote_path: "/"

`), []byte(`
name: "Tor router"
description: "Navigate through tor router securely"
icon_url: "https://upload.wikimedia.org/wikipedia/commons/7/73/Tor_logo-1.png"
remote_path: "/"
`), []byte(`
name: "Mantika VPN"
description: "Connect to company VPN"
icon_url: "http://soporte.fen.uchile.cl/mw/images/1/18/Vpnfen.png"
remote_path: "/"
`), []byte(`
name: "Metrics"
description: "Grafana dashboards"
icon_url: "https://demo.lightbend.com/grafana/public/img/grafana_icon.svg"
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
