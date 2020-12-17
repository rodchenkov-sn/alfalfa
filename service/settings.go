package service

import (
	"gopkg.in/ini.v1"
)

type CollectionPath struct {
	Database   string
	Collection string
}

type RepositorySettings struct {
	Uri               string
	UsersPath         CollectionPath
	MeasurementsPath  CollectionPath
}

type ServerSettings struct {
	RS         RepositorySettings
	PrivateKey string
}

func ReadSettings(file string) ServerSettings {
	settings, err := ini.Load(file)
	if err != nil {
		panic(err)
	}
	rs := RepositorySettings{
		Uri: settings.Section("").Key("db_link").String(),
		UsersPath: CollectionPath{
			Database:   settings.Section("users").Key("db_name").String(),
			Collection: settings.Section("users").Key("collection_name").String(),
		},
		MeasurementsPath: CollectionPath{
			Database:   settings.Section("measurements").Key("db_name").String(),
			Collection: settings.Section("measurements").Key("collection_name").String(),
		},
	}
	privateKey := settings.Section("").Key("private_key").String()
	return ServerSettings{
		RS: rs,
		PrivateKey: privateKey,
	}
}
