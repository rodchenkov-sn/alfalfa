package service

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
