package repos

import (
	"github.com/streamdp/ccd/db"
)

type SessionRepo struct {
	db db.Database
}

func NewSessionRepo(db db.Database) (db.Session, error) {
	return &SessionRepo{
		db: db,
	}, nil
}

func (sr *SessionRepo) UpdateTask(n string, i int64) (err error) {
	if _, err = sr.db.UpdateTask(n, i); err != nil {
		return
	}
	return
}

func (sr *SessionRepo) GetSession() (map[string]int64, error) {
	return sr.db.GetSession()
}

func (sr *SessionRepo) AddTask(n string, i int64) (err error) {
	if _, err = sr.db.AddTask(n, i); err != nil {
		return
	}
	return
}

func (sr *SessionRepo) RemoveTask(n string) (err error) {
	if _, err = sr.db.RemoveTask(n); err != nil {
		return
	}
	return
}
