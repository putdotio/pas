package sqlinserter

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/putdotio/pas/internal/property"
)

type Inserter interface {
	InsertSQL(def property.Types, timestamp time.Time) (sql string, values []interface{}, err error)
	CreateTableSQL(property.Types) (string, error)
	ExistingColumns(*sql.DB) (map[string]struct{}, error)
	AlterTableSQL(existingColumns map[string]struct{}, def property.Types) (string, error)
}

func Insert(i Inserter, t time.Time, db *sql.DB, def property.Types) error {
	sql, values, err := i.InsertSQL(def, t)
	if err != nil {
		return err
	}
	_, err = db.Exec(sql, values...)
	if merr, ok := err.(*mysql.MySQLError); ok {
		if merr.Number == 1146 { // table doesn't exist
			ctsql, err2 := i.CreateTableSQL(def)
			if err2 != nil {
				return err2
			}
			_, err = db.Exec(ctsql)
			if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1050 { // table already exists
				err = nil
			}
			if err != nil {
				return err
			}
			_, err = db.Exec(sql, values...)
		} else if merr.Number == 1054 { // unknown column
			cols, err2 := i.ExistingColumns(db)
			if err2 != nil {
				return err2
			}
			atsql, err2 := i.AlterTableSQL(cols, def)
			if err2 != nil {
				return err2
			}
			_, err = db.Exec(atsql)
			if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1060 { // duplicate column name
				err = nil
			}
			if err != nil {
				return err
			}
			_, err = db.Exec(sql, values...)
		}
	}
	return err
}
