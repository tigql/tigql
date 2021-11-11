package drivers

import (
	"github.com/feitian124/goapi/database/schema"
)

type Dialect int8

const (
	MYSQL = iota
	SQLITE
)

// Driver is the common interface for database drivers
type Driver interface {
	Analyze(*schema.Schema) error
	Info() (*schema.Driver, error)
	//Tables(schema string) ([]schema.Table, error)
	//Table(schema string, table string) (*schema.Table, error)
	//Columns(table string) ([]schema.Column, error)
}

// Option is the type for change Config.
type Option func(Driver) error
