package pas

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Analytics struct {
	db *sql.DB
}

func NewAnalytics(db *sql.DB) *Analytics {
	return &Analytics{
		db: db,
	}
}

func (p *Analytics) InsertEvents(events []Event) (n int, err error) {
	now := time.Now().UTC()
	for _, e := range events {
		i := eventInserter{e}
		err = p.insert(i, now)
		if err != nil {
			break
		}
		n++
	}
	return
}

func (p *Analytics) UpdateUsers(users []User) (n int, err error) {
	now := time.Now().UTC()
	for _, u := range users {
		i := userInserter{u}
		err = p.insert(i, now)
		if err != nil {
			break
		}
		n++
	}
	return
}

func (p *Analytics) insert(i inserter, t time.Time) error {
	sql, values := i.InsertSQL(t)
	_, err := p.db.Exec(sql, values...)
	if merr, ok := err.(*mysql.MySQLError); ok {
		if merr.Number == 1146 { // table doesn't exist
			_, err = p.db.Exec(i.CreateTableSQL())
			if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1050 { // table already exists
				err = nil
			}
			if err != nil {
				return err
			}
			_, err = p.db.Exec(sql, values...)
		} else if merr.Number == 1054 { // unknown column
			cols, err2 := i.ExistingColumns(p.db)
			if err2 != nil {
				return err2
			}
			_, err = p.db.Exec(i.AlterTableSQL(cols))
			if merr, ok = err.(*mysql.MySQLError); ok && merr.Number == 1060 { // duplicate column name
				err = nil
			}
			if err != nil {
				return err
			}
			_, err = p.db.Exec(sql, values...)
		}
	}
	return err
}

func (p *Analytics) Health() error {
	return p.db.Ping()
}
