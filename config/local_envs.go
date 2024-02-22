package config

import (
	"os"
	"strings"
	"sync"
)

var localEnvs *envs

func init() {
	localEnvs = newEnvs()
}

type envs struct {
	e map[string]string
	m *sync.RWMutex
}

func newEnvs() *envs {
	e := &envs{
		m: new(sync.RWMutex),
	}
	return e.init()
}

func (e *envs) init() *envs {
	var localEnvsArr = os.Environ()
	e.e = make(map[string]string, len(localEnvsArr))
	e.m.Lock()
	defer e.m.Unlock()
	for _, env := range localEnvsArr {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			e.e[pair[0]] = pair[1]
		}
	}
	return e
}

func (e *envs) get(name string) string {
	e.m.RLock()
	defer e.m.RUnlock()
	if env, ok := e.e[name]; ok {
		return env
	}
	return ""
}
