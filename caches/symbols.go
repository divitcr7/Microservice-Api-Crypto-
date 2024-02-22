package caches

import (
	"strings"
	"sync"
)

type SymbolCache struct {
	c map[string]struct{}
	m *sync.RWMutex
}

func NewSymbolCache() *SymbolCache {
	return &SymbolCache{
		c: make(map[string]struct{}),
		m: new(sync.RWMutex),
	}
}

func (sc *SymbolCache) Add(s string) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.c[strings.ToUpper(s)] = struct{}{}
}

func (sc *SymbolCache) Remove(s string) {
	s = strings.ToUpper(s)
	sc.m.Lock()
	defer sc.m.Unlock()
	if _, ok := sc.c[s]; ok {
		delete(sc.c, s)
	}
}

func (sc *SymbolCache) IsPresent(s string) bool {
	sc.m.RLock()
	defer sc.m.RUnlock()
	if _, ok := sc.c[strings.ToUpper(s)]; ok {
		return true
	}
	return false
}
