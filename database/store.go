package database

type DataStore interface {
	DumpDeleted() error
}
