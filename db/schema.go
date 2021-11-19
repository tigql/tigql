package db

// Schema is a named collection of tables. A schema can also contain views, indexes, sequences, data types, operators,
// and functions. Schemas are analogous to directories at the operating system level, except schemas cannot be nested.
// For mysql, schema is equal to database. For postgres, a database could have multi schemas.
type Schema struct {
	Name string `json:"name"`
}

type ISchema interface {
	Tables(condition string, params ...string) ([]*TableInfo, error)
	Table(table string) (*Table, error)
}

func (s *Schema) Tables(condition string, params ...string) ([]*TableInfo, error) {
	var tis []*TableInfo
	return tis, nil
}

func (s *Schema) Table(tableName string) (*Table, error) {
	t := &Table{}
	return t, nil
}