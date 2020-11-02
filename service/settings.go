package service

type CollectionSettings struct {
	Database   string
	Collection string
}

type RepositorySettings struct {
	Uri                  string
	UsersSettings        CollectionSettings
	MeasurementsSettings CollectionSettings
}
