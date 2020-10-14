package userinserter

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"pas/internal/inserter"
	"pas/internal/property"
	"pas/internal/user"
)

type UserInserter struct {
	User user.User
}

var _ inserter.Inserter = (*UserInserter)(nil)

func (i UserInserter) InsertSQL(def property.Types, t time.Time) (string, []interface{}, error) {
	u := i.User
	var sb strings.Builder
	sb.WriteString("insert into user")
	sb.WriteString("(id, timestamp")
	values := make([]interface{}, 0, 2*len(u.Properties)+3)
	values = append(values, string(u.ID))
	values = append(values, t)
	for pname, pval := range u.Properties {
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
	sb.WriteString(") values (?, ?")
	for range u.Properties {
		sb.WriteString(",?")
	}
	sb.WriteString(") on duplicate key update timestamp = ?")
	values = append(values, t)
	for pname, pval := range u.Properties {
		sb.WriteRune(',')
		sb.WriteString(string(pname))
		sb.WriteString("=?")
		if pval != nil {
			val, err := def[pname].ConvertValue(pval)
			if err != nil {
				return "", nil, err
			}
			values = append(values, val)
		} else {
			values = append(values, nil)
		}
	}
	return sb.String(), values, nil
}

func (i UserInserter) CreateTableSQL(def property.Types) (string, error) {
	u := i.User
	var sb strings.Builder
	sb.WriteString("create table user")
	sb.WriteString("(id varchar(255) not null, timestamp datetime not null")
	for pname := range u.Properties {
		ptype, ok := def[pname]
		if !ok {
			return "", errors.New("unknown property: " + string(pname))
		}
		sb.WriteRune(',')
		sb.WriteString(string(pname))
		sb.WriteRune(' ')
		sb.WriteString(ptype.ColumnType())
	}
	sb.WriteString(", primary key (id), index idx_userid (id), index idx_timestamp(timestamp))")
	return sb.String(), nil
}

func (i UserInserter) ExistingColumns(db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.Query("select column_name from information_schema.columns where table_name = 'user' and column_name != 'id'")
	if err != nil {
		return nil, err
	}
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

func (i UserInserter) AlterTableSQL(existingColumns map[string]struct{}, def property.Types) (string, error) {
	u := i.User
	var sb strings.Builder
	sb.WriteString("alter table user ")
	for pname := range u.Properties {
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
