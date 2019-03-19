package analytics

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/putdotio/pas/internal/event"
	"github.com/putdotio/pas/internal/inserter"
	"github.com/putdotio/pas/internal/inserter/eventinserter"
	"github.com/putdotio/pas/internal/inserter/userinserter"
	"github.com/putdotio/pas/internal/property"
	"github.com/putdotio/pas/internal/user"
)

type Analytics struct {
	db     *sql.DB
	schema schema
	secret []byte
}

type schema struct {
	user   property.Types
	events map[event.Name]property.Types
}

func New(db *sql.DB, secret string, user property.Types, events map[event.Name]property.Types) *Analytics {
	return &Analytics{
		db:     db,
		schema: schema{user, events},
		secret: []byte(secret),
	}
}

func (p *Analytics) InsertEvents(events []event.Event) (n int, err error) {
	hash := hmac.New(sha256.New, p.secret)
	now := time.Now().UTC()
	for _, e := range events {
		if e.UserHash == nil {
			e.IsAnonymous = true
		} else {
			hash.Reset()
			_, _ = hash.Write([]byte(e.UserID))
			if hex.EncodeToString(hash.Sum(nil)) != *e.UserHash {
				err = errors.New("invalid user hash: " + *e.UserHash)
				return
			}
		}
		err = p.insertEvent(e, now)
		if err != nil {
			return
		}
		n++
	}
	return
}

func (p *Analytics) UpdateUsers(users []user.User) (n int, err error) {
	hash := hmac.New(sha256.New, p.secret)
	now := time.Now().UTC()
	for _, u := range users {
		hash.Reset()
		_, _ = hash.Write([]byte(u.ID))
		if hex.EncodeToString(hash.Sum(nil)) != u.Hash {
			err = errors.New("invalid hash: " + u.Hash)
			return
		}
		err = p.insertUser(u, now)
		if err != nil {
			return
		}
		n++
	}
	return
}

func (p *Analytics) insertEvent(e event.Event, t time.Time) error {
	def, ok := p.schema.events[e.Name]
	if !ok {
		return errors.New("unknown event name: " + string(e.Name))
	}
	i := eventinserter.EventInserter{Event: e}
	return inserter.Insert(i, t, p.db, def)
}

func (p *Analytics) insertUser(u user.User, t time.Time) error {
	i := userinserter.UserInserter{User: u}
	return inserter.Insert(i, t, p.db, p.schema.user)
}

func (p *Analytics) Health() error {
	return p.db.Ping()
}

func (p *Analytics) Alias(previousID, userID user.ID, userHash string) error {
	hash := hmac.New(sha256.New, p.secret)
	_, _ = hash.Write([]byte(userID))
	if hex.EncodeToString(hash.Sum(nil)) != userHash {
		return errors.New("invalid hash: " + userHash)
	}
	sql := "SELECT table_name FROM INFORMATION_SCHEMA.tables WHERE table_schema = (SELECT DATABASE()) AND table_name != 'user'"
	rows, err := p.db.Query(sql)
	if err != nil {
		return err
	}
	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return err
		}
		tables = append(tables, table)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	for _, table := range tables {
		sql = fmt.Sprintf("update %s set user_id = ?, is_anonymous=0 where user_id = ?", table)
		_, err = tx.Exec(sql, userID, previousID)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
