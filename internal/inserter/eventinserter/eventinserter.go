package eventinserter

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/putdotio/pas/internal/event"
	"github.com/putdotio/pas/internal/inserter"
	"github.com/putdotio/pas/internal/property"
)

type EventInserter struct {
	Event event.Event
}

var _ inserter.Inserter = (*EventInserter)(nil)

func (i EventInserter) InsertSQL(def property.Types, t time.Time) (string, []interface{}, error) {
	e := i.Event
	var sb strings.Builder
	sb.WriteString("insert into ")
	sb.WriteString(string(e.Name))
	sb.WriteString("(user_id, timestamp, is_anonymous")
	values := make([]interface{}, 0, len(e.Properties)+3)
	values = append(values, string(e.UserID))
	if e.Timestamp != nil {
		values = append(values, e.Timestamp)
	} else {
		values = append(values, t)
	}
	values = append(values, e.IsAnonymous)
	for pname, pval := range e.Properties {
		ptype, ok := def[pname]
		if !ok {
			return "", nil, errors.New("unknown property: " + string(pname))
		}
		sb.WriteRune(',')
		sb.WriteString(string(pname))
		if pval != nil {
			val, err := ptype.ConvertValue(pval)
			if err != nil {
				return "", nil, errors.New("cannot read property (" + string(pname) + "): " + err.Error())
			}
			values = append(values, val)
		} else {
			values = append(values, nil)
		}
	}
	sb.WriteString(") values (?, ?, ?")
	for range e.Properties {
		sb.WriteString(",?")
	}
	sb.WriteRune(')')
	return sb.String(), values, nil
}

func (i EventInserter) CreateTableSQL(def property.Types) (string, error) {
	e := i.Event
	var sb strings.Builder
	sb.WriteString("create table ")
	sb.WriteString(string(e.Name))
	sb.WriteString("(user_id varchar(255) not null, timestamp datetime not null, is_anonymous tinyint(1) not null")
	for pname := range e.Properties {
		ptype, ok := def[pname]
		if !ok {
			return "", errors.New("unknown property: " + string(pname))
		}
		sb.WriteRune(',')
		sb.WriteString(string(pname))
		sb.WriteRune(' ')
		sb.WriteString(ptype.ColumnType())
	}
	sb.WriteString(", index idx_userid_timestamp (user_id, timestamp), index idx_timestamp(timestamp))")
	return sb.String(), nil
}

func (i EventInserter) ExistingColumns(db *sql.DB) (map[string]struct{}, error) {
	table := string(i.Event.Name)
	rows, err := db.Query("select column_name from information_schema.columns where table_name = ? and column_name not in ('user_id', 'timestamp', 'is_anonymous')", table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	existingColumns := make(map[string]struct{})
	for rows.Next() {
		var col string
		err = rows.Scan(&col)
		if err != nil {
			return nil, err
		}
		existingColumns[col] = struct{}{}
	}
	return existingColumns, rows.Err()
}

func (i EventInserter) AlterTableSQL(existingColumns map[string]struct{}, def property.Types) (string, error) {
	e := i.Event
	var sb strings.Builder
	sb.WriteString("alter table ")
	sb.WriteString(string(e.Name))
	for pname := range e.Properties {
		ptype, ok := def[pname]
		if !ok {
			return "", errors.New("unknown property: " + string(pname))
		}
		_, ok = existingColumns[string(pname)]
		if !ok {
			sb.WriteString(" add column ")
			sb.WriteString(string(pname))
			sb.WriteRune(' ')
			sb.WriteString(ptype.ColumnType())
			sb.WriteRune(',')
		}
	}
	return strings.TrimRight(sb.String(), ","), nil
}
