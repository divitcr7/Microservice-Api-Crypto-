package postgres

import (
	"database/sql"
	"errors"
	"strings"
)

func (d *Db) AddTask(n string, i int64) (result sql.Result, err error) {
	if n == "" {
		return nil, errors.New("cant insert empty task name")
	}
	return d.Exec(
		`insert into session (task_name,interval) values ($1,$2) on conflict do nothing;`, strings.ToUpper(n), i,
	)
}

func (d *Db) UpdateTask(n string, i int64) (result sql.Result, err error) {
	if n == "" {
		return nil, errors.New("empty task name")
	}
	return d.Exec(`update session set interval=$2 where task_name=$1;`, strings.ToUpper(n), i)
}

func (d *Db) RemoveTask(n string) (result sql.Result, err error) {
	if n == "" {
		return nil, errors.New("empty symbol")
	}
	return d.Exec(`delete from session where task_name=$1;`, strings.ToUpper(n))
}

func (d *Db) GetSession() (tasks map[string]int64, err error) {
	rows, err := d.Query(`select task_name,interval from session`)
	if err != nil {
		return nil, err
	}
	tasks = make(map[string]int64)
	for rows.Next() {
		var (
			n string
			i int64
		)
		if err = rows.Scan(&n, &i); err != nil {
			return nil, err
		}
		tasks[n] = i
	}
	return
}
