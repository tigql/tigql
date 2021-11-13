package datasource

import (
	"strings"

	"github.com/feitian124/goapi/config"
	"github.com/feitian124/goapi/database/drivers"
	"github.com/feitian124/goapi/database/drivers/mariadb"
	"github.com/feitian124/goapi/database/drivers/mysql"
	"github.com/feitian124/goapi/database/drivers/postgres"
	"github.com/feitian124/goapi/database/drivers/sqlite"
	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

func Analyze(dsn config.DSN) (*schema.Schema, error) {
	url := dsn.URL
	s := &schema.Schema{}
	u, err := dburl.Parse(url)
	if err != nil {
		return s, errors.WithStack(err)
	}
	parts := strings.Split(u.Short(), "/")
	if len(parts) < 2 {
		return s, errors.Errorf("invalid DSN: parse %s -> %#v", url, u)
	}

	var opts []drivers.Option
	if u.Driver == "mysql" {
		values := u.Query()
		for k := range values {
			if k == "show_auto_increment" {
				opts = append(opts, mysql.ShowAutoIcrrement())
				values.Del(k)
			}
			if k == "hide_auto_increment" {
				opts = append(opts, mysql.HideAutoIcrrement())
				values.Del(k)
			}
		}
		u.RawQuery = values.Encode()
		url = u.String()
	}

	db, err := dburl.Open(url)
	defer db.Close()
	if err != nil {
		return s, errors.WithStack(err)
	}
	if err := db.Ping(); err != nil {
		return s, errors.WithStack(err)
	}

	var driver drivers.Driver

	switch u.Driver {
	case "mysql":
		s.Name = parts[1]
		if u.Scheme == "maria" || u.Scheme == "mariadb" {
			driver, err = mariadb.New(db, opts...)
		} else {
			driver, err = mysql.New(db, opts...)
		}
		if err != nil {
			return s, err
		}
	case "postgres":
		s.Name = parts[1]
		driver = postgres.New(db)
	case "sqlite3":
		s.Name = parts[len(parts)-1]
		driver = sqlite.New(db)
	default:
		return s, errors.Errorf("unsupported driver '%s'", u.Driver)
	}
	err = driver.Analyze(s)
	if err != nil {
		return s, err
	}
	return s, nil
}
