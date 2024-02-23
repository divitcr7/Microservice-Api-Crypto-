package repos

import (
	"github.com/streamdp/ccd/caches"
	"github.com/streamdp/ccd/db"
)

type SymbolRepo struct {
	db db.Database
	c  *caches.SymbolCache
}

func NewSymbolRepository(db db.Database) *SymbolRepo {
	return &SymbolRepo{
		c:  caches.NewSymbolCache(),
		db: db,
	}
}

func (sc *SymbolRepo) Update(s, u string) (err error) {
	if _, err = sc.db.UpdateSymbol(s, u); err != nil {
		return
	}
	sc.c.Add(s)
	return
}

func (sc *SymbolRepo) Load() error {
	s, err := sc.db.Symbols()
	if err != nil {
		return err
	}
	for i := range s {
		sc.c.Add(s[i].Symbol)
	}
	return nil
}

func (sc *SymbolRepo) Add(s, u string) (err error) {
	if _, err = sc.db.AddSymbol(s, u); err != nil {
		return
	}
	sc.c.Add(s)
	return
}

func (sc *SymbolRepo) Remove(s string) (err error) {
	if _, err = sc.db.RemoveSymbol(s); err != nil {
		return
	}
	sc.c.Remove(s)
	return
}

func (sc *SymbolRepo) IsPresent(s string) bool {
	return sc.c.IsPresent(s)
}
