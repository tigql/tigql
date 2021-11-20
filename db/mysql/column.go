package mysql

import (
	"database/sql"
	"fmt"

	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
)

const supportGeneratedColumnSQL = `
	SELECT column_name, column_default, is_nullable, column_type, column_comment, extra, generation_expression
	FROM information_schema.columns
	WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position
`

const columnSQL = `
	SELECT column_name, column_default, is_nullable, column_type, column_comment, extra
	FROM information_schema.columns
	WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position
`

func (d *DB) Columns(tableName string) ([]*schema.Column, error) {
	columnStmt := supportGeneratedColumnSQL
	if !d.supportGeneratedColumn {
		columnStmt = columnSQL
	}
	columnRows, err := d.db.Query(columnStmt, d.Schema.Name, tableName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer columnRows.Close()
	var columns []*schema.Column
	for columnRows.Next() {
		var (
			columnName     string
			columnDefault  sql.NullString
			isNullable     string
			columnType     string
			columnComment  sql.NullString
			extra          sql.NullString
			generationExpr sql.NullString
		)
		if d.supportGeneratedColumn {
			err = columnRows.Scan(&columnName, &columnDefault, &isNullable, &columnType, &columnComment, &extra, &generationExpr)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			err = columnRows.Scan(&columnName, &columnDefault, &isNullable, &columnType, &columnComment, &extra)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
		extraDef := extra.String
		if generationExpr.String != "" {
			switch extraDef {
			case "VIRTUAL GENERATED":
				extraDef = fmt.Sprintf("GENERATED ALWAYS AS %s VIRTUAL", generationExpr.String)
			case "STORED GENERATED":
				extraDef = fmt.Sprintf("GENERATED ALWAYS AS %s STORED", generationExpr.String)
			default:
				extraDef = fmt.Sprintf("%s:%s", extraDef, generationExpr.String)
			}
		}
		column := &schema.Column{
			Name:     columnName,
			Type:     columnType,
			Nullable: convertColumnNullable(isNullable),
			Default:  columnDefault,
			Comment:  columnComment.String,
			ExtraDef: extraDef,
		}

		columns = append(columns, column)
	}
	return columns, nil
}

func convertColumnNullable(str string) bool {
	return str != "NO"
}