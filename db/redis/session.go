package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/streamdp/ccd/config"
)

const sessionName = "lastSession"

type KeysStore struct {
	c *redis.Client
}

// NewRedisKeysStore initialize new redis session store
func NewRedisKeysStore() (*KeysStore, error) {
	opt, err := getRedisOptions()
	if err != nil {
		return nil, fmt.Errorf("filed to parse redis os environment variables: %w", err)
	}
	client := redis.NewClient(opt)
	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}
	return &KeysStore{
		c: client,
	}, nil
}

func getSeparatedOptions() (*redis.Options, error) {
	var (
		host     = "127.0.0.1"
		port     = 6379
		password = ""
		db       = 0
	)
	if h := config.GetEnv("REDIS_HOSTNAME"); h != "" {
		host = h
	}
	if p := config.GetEnv("REDIS_PORT"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		port = n
	}
	if pass := config.GetEnv("REDIS_PASSWORD"); pass != "" {
		password = pass
	}
	if d := config.GetEnv("REDIS_DB"); d != "" {
		n, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}
		db = n
	}
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	}, nil
}

func getRedisOptions() (*redis.Options, error) {
	if redisUrl := config.GetEnv("REDIS_URL"); redisUrl != "" {
		return redis.ParseURL(config.GetEnv("REDIS_URL"))
	}
	return getSeparatedOptions()
}

// GetSession get previously saved session
func (s *KeysStore) GetSession() (map[string]int64, error) {
	if s == nil {
		return nil, nil
	}
	session := make(map[string]int64)
	for k, v := range s.c.HGetAll(sessionName).Val() {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		session[k] = i
	}
	return session, nil
}

// AddTask add a new task or update an already saved task in the current session
func (s *KeysStore) AddTask(n string, i int64) (err error) {
	if s == nil {
		return
	}
	s.c.HSet(sessionName, n, strconv.FormatInt(i, 10))
	return
}

// UpdateTask add a new task or update an already saved task in the current session
func (s *KeysStore) UpdateTask(n string, i int64) (err error) {
	return s.AddTask(n, i)
}

// RemoveTask remove a task from the current session
func (s *KeysStore) RemoveTask(n string) (err error) {
	if s == nil {
		return
	}
	s.c.HDel(sessionName, n)
	return
}
