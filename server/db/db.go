package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mobyos/admin-app/server/types"
	"github.com/twinj/uuid"
)

func GetInstallations() ([]*types.Installation, error) {
	db, err := sql.Open("sqlite3", "./db/mobyos.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select id, application_id from installation")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []*types.Installation{}
	for rows.Next() {
		var id string
		var applicationId string
		err = rows.Scan(&id, &applicationId)
		if err != nil {
			return nil, err
		}
		apps = append(apps, &types.Installation{Id: id, ApplicationId: applicationId})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func GetInstallation(appId string) (*types.Installation, error) {
	db, err := sql.Open("sqlite3", "./db/mobyos.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow("select id, application_id, descriptor from installation where id = ?", appId)

	app := &types.Installation{}
	err = row.Scan(&app.Id, &app.ApplicationId, &app.Descriptor)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return app, nil
}

func CreateApplication(appDesc types.AppDescriptor) error {
	db, err := sql.Open("sqlite3", "./db/mobyos.db")
	if err != nil {
		return err
	}
	defer db.Close()

	desc, err := appDesc.GetBytes()
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into installation values (?, ?, ?)", uuid.NewV4().String(), appDesc.Name, desc)
	if err != nil {
		return err
	}

	return nil
}

func DeleteApplication(appId string) error {
	db, err := sql.Open("sqlite3", "./db/mobyos.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("delete from installation where id = (?)", appId)
	if err != nil {
		return err
	}

	return nil
}
