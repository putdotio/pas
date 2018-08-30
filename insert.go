package main

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

func runInserter(i Inserter, t time.Time) error {
	sql, values := i.InsertSQL(t)
	_, err := db.Exec(sql, values...)
	if merr, ok := err.(*mysql.MySQLError); ok {
		if merr.Number == 1146 { // table doesn't exist
			_, err = db.Exec(i.CreateTableSQL())
			if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1050 { // table already exists
				err = nil
			}
			if err != nil {
				return err
			}
			_, err = db.Exec(sql, values...)
		} else if merr.Number == 1054 { // unknown column
			cols, err2 := i.ExistingColumns()
			if err2 != nil {
				return err2
			}
			_, err = db.Exec(i.AlterTableSQL(cols))
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

type Inserter interface {
	InsertSQL(timestamp time.Time) (sql string, values []interface{})
	CreateTableSQL() string
	ExistingColumns() (map[string]struct{}, error)
	AlterTableSQL(existingColumns map[string]struct{}) string
}

type EventInserter struct {
	e Event
}

func (i EventInserter) InsertSQL(t time.Time) (sql string, values []interface{}) {
	return insertEvent(i.e, t)
}

func (i EventInserter) CreateTableSQL() string {
	return createEventTable(i.e)
}

func (i EventInserter) ExistingColumns() (map[string]struct{}, error) {
	return existingEventColumns(string(i.e.Name))
}

func (i EventInserter) AlterTableSQL(existingColumns map[string]struct{}) string {
	return alterEventTable(i.e, existingColumns)
}

type UserInserter struct {
	u User
}

func (i UserInserter) InsertSQL(t time.Time) (sql string, values []interface{}) {
	return insertUser(i.u, t)
}

func (i UserInserter) CreateTableSQL() string {
	return createUserTable(i.u)
}

func (i UserInserter) ExistingColumns() (map[string]struct{}, error) {
	return existingUserColumns()
}

func (i UserInserter) AlterTableSQL(existingColumns map[string]struct{}) string {
	return alterUserTable(i.u, existingColumns)
}
