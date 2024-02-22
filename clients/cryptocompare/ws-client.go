package cryptocompare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/streamdp/ccd/domain"
	"nhooyr.io/websocket"

	"github.com/streamdp/ccd/clients"
)

const wssUrl = "wss://streamer.cryptocompare.com/v2"

type cryptoCompareWs struct {
	ctx        context.Context
	l          *log.Logger
	conn       *websocket.Conn
	apiKey     string
	subscribes domain.Subscribes
	subMu      sync.RWMutex
}

func InitWs(pipe chan *domain.Data, l *log.Logger) (_ clients.WsClient, err error) {
	var apiKey string
	if apiKey, err = getApiKey(); err != nil {
		return nil, err
	}
	h := &cryptoCompareWs{
		ctx:        context.Background(),
		l:          l,
		apiKey:     apiKey,
		subscribes: domain.Subscribes{},
	}
	if err = h.reconnect(); err != nil {
		return nil, err
	}
	h.handleWsMessages(pipe)
	return h, nil
}

func (c *cryptoCompareWs) reconnect() (err error) {
	if c.conn != nil {
		if err = c.conn.Close(websocket.StatusNormalClosure, ""); err != nil {
			c.l.Println(err)
			// reducing logs and CPU load when API key expired
			time.Sleep(10 * time.Second)
		}
	}
	var u *url.URL
	if u, err = c.buildURL(); err != nil {
		return
	}
	c.conn, _, err = websocket.Dial(c.ctx, u.String(), nil)
	return
}

func (c *cryptoCompareWs) buildURL() (u *url.URL, err error) {
	if u, err = url.Parse(wssUrl); err != nil {
		return
	}
	query := u.Query()
	query.Set("api_key", c.apiKey)
	u.RawQuery = query.Encode()
	return
}

func (c *cryptoCompareWs) resubscribe() (err error) {
	c.subMu.RLock()
	defer c.subMu.RUnlock()
	for k := range c.subscribes {
		if err = c.sendSubscribeMsg(k); err != nil {
			return
		}
	}
	return
}

func (c *cryptoCompareWs) handleWssError(err error) error {
	c.l.Println(err)
	for {
		select {
		case <-time.After(time.Minute):
			return errors.New("reconnect failed")
		default:
			if err = c.reconnect(); err != nil {
				time.Sleep(time.Second)
				continue
			}
			if err = c.resubscribe(); err != nil {
				time.Sleep(time.Second)
				continue
			}
			return nil
		}
	}
}

func (c *cryptoCompareWs) handleWsMessages(pipe chan *domain.Data) {
	go func() {
		defer func(conn *websocket.Conn, code websocket.StatusCode, reason string) {
			if err := conn.Close(code, reason); err != nil {
				c.l.Println(err)
			}
		}(c.conn, websocket.StatusNormalClosure, "")
		var (
			// CryptoCompare automatically send a heartbeat message per socket every 30 seconds,
			// if we miss two heartbeats it means our connection might be stale. So let's start with 2 and
			// add 1 every time the server sends a heartbeat and subtract 1 by ticker every 30 seconds.
			hb   = 2
			tick = time.NewTicker(30 * time.Second)
		)
		defer tick.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-tick.C:
				if hb <= 0 {
					if err := c.handleWssError(errors.New("heartbeat loss")); err != nil {
						c.l.Println(err)
						return
					}
				}
				hb--
			default:
				var (
					body []byte
					err  error
				)
				if _, body, err = c.conn.Read(c.ctx); err != nil {
					if err = c.handleWssError(err); err != nil {
						c.l.Println(err)
						return
					}
					continue
				}
				data := &cryptoCompareWsData{}
				if err = json.Unmarshal(body, data); err != nil {
					c.l.Println(err)
					continue
				}
				switch data.Type {
				case "999":
					hb++
				case "5":
					pipe <- convertCryptoCompareWsDataToDomain(data)
				}
			}
		}
	}()
}

func buildChannelName(from, to string) string {
	return fmt.Sprintf("5~CCCAGG~%s~%s", strings.ToUpper(from), strings.ToUpper(to))
}

func (c *cryptoCompareWs) Unsubscribe(from, to string) (err error) {
	c.subMu.Lock()
	defer c.subMu.Unlock()
	var ch = buildChannelName(from, to)
	if _, ok := c.subscribes[ch]; ok {
		if err = c.sendUnsubscribeMsg(ch); err != nil {
			return
		}
		delete(c.subscribes, ch)
	}
	return
}

func (c *cryptoCompareWs) sendUnsubscribeMsg(ch string) error {
	return c.conn.Write(c.ctx, websocket.MessageText, []byte(
		fmt.Sprintf("{\"action\":\"SubRemove\",\"subs\":[\"%s\"]}", ch)),
	)
}

func (c *cryptoCompareWs) Subscribe(from, to string) (err error) {
	c.subMu.Lock()
	defer c.subMu.Unlock()
	var ch = buildChannelName(from, to)
	if err = c.sendSubscribeMsg(ch); err != nil {
		return
	}
	c.subscribes[ch] = domain.NewSubscribe(from, to, 0)
	return
}

func (c *cryptoCompareWs) sendSubscribeMsg(ch string) error {
	return c.conn.Write(c.ctx, websocket.MessageText, []byte(
		fmt.Sprintf("{\"action\":\"SubAdd\",\"subs\":[\"%s\"]}", ch)),
	)
}

func (c *cryptoCompareWs) ListSubscribes() domain.Subscribes {
	s := make(domain.Subscribes, len(c.subscribes))
	c.subMu.RLock()
	defer c.subMu.RUnlock()
	for k, v := range c.subscribes {
		s[k] = v
	}
	return s
}

func convertCryptoCompareWsDataToDomain(d *cryptoCompareWsData) *domain.Data {
	if d == nil {
		return nil
	}
	b, _ := json.Marshal(&domain.Raw{
		FromSymbol:     d.FromSymbol,
		ToSymbol:       d.ToSymbol,
		Open24Hour:     d.Open24Hour,
		Volume24Hour:   d.Volume24Hour,
		Volume24HourTo: d.Volume24HourTo,
		High24Hour:     d.High24Hour,
		Price:          d.Price,
		LastUpdate:     d.LastUpdate,
		Supply:         d.CurrentSupply,
		MktCap:         d.CurrentSupplyMktCap,
	})
	return &domain.Data{
		FromSymbol:     d.FromSymbol,
		ToSymbol:       d.ToSymbol,
		Open24Hour:     d.Open24Hour,
		Volume24Hour:   d.Volume24Hour,
		Low24Hour:      d.Low24Hour,
		High24Hour:     d.High24Hour,
		Price:          d.Price,
		Supply:         d.CurrentSupply,
		MktCap:         d.CurrentSupplyMktCap,
		LastUpdate:     d.LastUpdate,
		DisplayDataRaw: string(b),
	}
}
