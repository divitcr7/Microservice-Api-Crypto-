package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/streamdp/ccd/clients"
	"github.com/streamdp/ccd/clients/cryptocompare"
	"github.com/streamdp/ccd/clients/huobi"
	"github.com/streamdp/ccd/config"
	"github.com/streamdp/ccd/db"
	"github.com/streamdp/ccd/db/redis"
	"github.com/streamdp/ccd/repos"
	"github.com/streamdp/ccd/router"
)

func main() {
	l := log.New(os.Stderr, "CCD:", log.LstdFlags)

	config.ParseFlags()
	gin.SetMode(config.RunMode)

	d, err := db.Connect(l)
	if err != nil {
		l.Fatalln(err)
	}

	s := newSessionStore(d, l)

	sr := repos.NewSymbolRepository(d)
	if err = sr.Load(); err != nil {
		l.Fatalln(err)
	}

	r, err := initRestClient()
	if err != nil {
		l.Fatalln(err)
	}

	w, err := initWsClient(d, l)
	if err != nil {
		l.Fatalln(err)
	}

	p := clients.NewPuller(r, l, s, d.DataPipe())
	if err = p.RestoreLastSession(); err != nil {
		l.Println(fmt.Errorf("error restoring last session: %w", err))
	}

	e := gin.Default()
	if err = router.InitRouter(e, d, l, sr, r, w, p); err != nil {
		l.Fatalln(err)
	}
	if err = e.Run(config.Port); err != nil {
		l.Fatalln(err)
	}
}

func initRestClient() (r clients.RestClient, err error) {
	switch config.DataProvider {
	case "huobi":
		return huobi.Init()
	default:
		return cryptocompare.Init()
	}
}

func initWsClient(d db.Database, l *log.Logger) (w clients.WsClient, err error) {
	switch config.DataProvider {
	case "huobi":
		return huobi.InitWs(d.DataPipe(), l)
	default:
		return cryptocompare.InitWs(d.DataPipe(), l)
	}
}

func newSessionStore(d db.Database, l *log.Logger) (s db.Session) {
	var err error
	switch config.SessionStore {
	case "redis":
		s, err = redis.NewRedisKeysStore()
	default:
		s, err = repos.NewSessionRepo(d)
	}
	if err != nil {
		l.Println(fmt.Errorf("failed to init session store: %w", err))
	}
	return
}
