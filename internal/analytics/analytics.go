package analytics

import (
	"database/sql"
	"errors"
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
}

type schema struct {
	user   property.Types
	events map[event.Name]property.Types
}

func New(db *sql.DB, user property.Types, events map[event.Name]property.Types) *Analytics {
	return &Analytics{
		db:     db,
		schema: schema{user, events},
	}
}

func (p *Analytics) InsertEvents(events []event.Event) (n int, err error) {
	now := time.Now().UTC()
	for _, e := range events {
		err = p.insertEvent(e, now)
		if err != nil {
			break
		}
		n++
	}
	return
}

func (p *Analytics) UpdateUsers(users []user.User) (n int, err error) {
	now := time.Now().UTC()
	for _, u := range users {
		err = p.insertUser(u, now)
		if err != nil {
			break
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
