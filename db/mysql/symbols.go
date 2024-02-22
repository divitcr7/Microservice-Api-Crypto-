package mysql

import (
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/streamdp/ccd/domain"
)

func (d *Db) AddSymbol(s string, u string) (result sql.Result, err error) {
	if s == "" {
		return nil, errors.New("cant insert empty symbol")
	}
	return d.Exec(`insert into symbols (symbol,unicode) values (?,?);`, strings.ToUpper(s), strings.ToUpper(u))
}

func (d *Db) UpdateSymbol(s, u string) (result sql.Result, err error) {
	if s == "" {
		return nil, errors.New("empty symbol")
	}
	return d.Exec(`update symbols set unicode=? where symbol=?;`, strings.ToUpper(u), strings.ToUpper(s))
}

func (d *Db) RemoveSymbol(s string) (result sql.Result, err error) {
	if s == "" {
		return nil, errors.New("empty symbol")
	}
	return d.Exec(`delete from symbols where symbol=?;`, strings.ToUpper(s))
}

func (d *Db) Symbols() (symbols []*domain.Symbol, err error) {
	rows, err := d.Query(`select symbol, unicode from symbols`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			s string
			u []byte
		)
		if err = rows.Scan(&s, &u); err != nil {
			return nil, err
		}
		r, _ := utf8.DecodeRune(u)
		symbols = append(symbols, &domain.Symbol{
			Symbol:  s,
			Unicode: r,
		})
	}
	return
}
