package huobi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/streamdp/ccd/domain"
	"nhooyr.io/websocket"

	"github.com/streamdp/ccd/clients"
)

const wssUrl = "wss://api.huobi.pro/ws"

type huobiWs struct {
	ctx        context.Context
	l          *log.Logger
	conn       *websocket.Conn
	subscribes domain.Subscribes
	subMu      sync.RWMutex
}

func InitWs(pipe chan *domain.Data, l *log.Logger) (clients.WsClient, error) {
	h := &huobiWs{
		ctx:        context.Background(),
		l:          l,
		subscribes: domain.Subscribes{},
	}
	if err := h.reconnect(); err != nil {
		return nil, err
	}
	h.handleWsMessages(pipe)
	return h, nil
}

func (h *huobiWs) reconnect() (err error) {
	if h.conn != nil {
		if err := h.conn.Close(websocket.StatusNormalClosure, ""); err != nil {
			h.l.Println(err)
			// reducing logs and CPU load when API key expired
			time.Sleep(10 * time.Second)
		}
	}
	h.conn, _, err = websocket.Dial(h.ctx, wssUrl, nil)
	return
}

func (h *huobiWs) resubscribe() (err error) {
	h.subMu.RLock()
	defer h.subMu.RUnlock()
	for k, v := range h.subscribes {
		if err = h.sendSubscribeMsg(k, v.Id()); err != nil {
			return
		}
	}
	return
}

func (h *huobiWs) handleWsError(err error) error {
	h.l.Println(err)
	for {
		select {
		case <-time.After(time.Minute):
			return errors.New("reconnect failed")
		default:
			if err = h.reconnect(); err != nil {
				time.Sleep(time.Second)
				continue
			}
			if err = h.resubscribe(); err != nil {
				time.Sleep(time.Second)
				continue
			}
			return nil
		}
	}
}

func (h *huobiWs) handleWsMessages(pipe chan *domain.Data) {
	go func() {
		defer func(conn *websocket.Conn, code websocket.StatusCode, reason string) {
			if err := conn.Close(code, reason); err != nil {
				h.l.Println(err)
			}
		}(h.conn, websocket.StatusNormalClosure, "")
		for {
			select {
			case <-h.ctx.Done():
				return
			default:
				var (
					r    io.Reader
					body []byte
					err  error
				)
				if _, r, err = h.conn.Reader(h.ctx); err != nil {
					if err = h.handleWsError(err); err != nil {
						h.l.Println(err)
						return
					}
					continue
				}
				if body, err = gzipDecompress(r); err != nil {
					h.l.Println(err)
					continue
				}
				if bytes.Contains(body, []byte("ping")) {
					if err = h.pingHandler(body); err != nil {
						if err = h.handleWsError(err); err != nil {
							h.l.Println(err)
							return
						}
					}
					continue
				}
				data := &huobiWsData{}
				if err = json.Unmarshal(body, data); err != nil {
					h.l.Println(err)
					continue
				}
				if data.Ch == "" {
					continue
				}
				from, to := h.pairFromChannelName(data.Ch)
				if from != "" && to != "" {
					pipe <- convertHuobiWsDataToDomain(from, to, data)
				}
			}
		}
	}()
}

func (h *huobiWs) pingHandler(m []byte) (err error) {
	m = bytes.Replace(m, []byte("ping"), []byte("pong"), -1)
	return h.conn.Write(h.ctx, websocket.MessageText, m)
}

func (h *huobiWs) pairFromChannelName(ch string) (from, to string) {
	h.subMu.RLock()
	defer h.subMu.RUnlock()
	if c, ok := h.subscribes[ch]; ok {
		return c.From, c.To
	}
	return
}

func buildChannelName(from, to string) string {
	if strings.ToLower(to) == "usd" {
		to = "usdt"
	}
	return fmt.Sprintf("market.%s.ticker", strings.ToLower(from+to))
}

func (h *huobiWs) Unsubscribe(from, to string) (err error) {
	h.subMu.Lock()
	defer h.subMu.Unlock()
	var ch = buildChannelName(from, to)
	if c, ok := h.subscribes[ch]; ok {
		if err = h.sendUnsubscribeMsg(ch, c.Id()); err != nil {
			return
		}
		delete(h.subscribes, ch)
	}
	return
}

func (h *huobiWs) sendUnsubscribeMsg(ch string, id int64) error {
	return h.conn.Write(h.ctx, websocket.MessageText, []byte(
		fmt.Sprintf("{\"unsub\": \"%s\", \"id\":\"%d\"}", ch, id)),
	)
}

func (h *huobiWs) Subscribe(from, to string) (err error) {
	h.subMu.Lock()
	defer h.subMu.Unlock()
	var (
		id = time.Now().UnixMilli()
		ch = buildChannelName(from, to)
	)
	if err = h.sendSubscribeMsg(ch, id); err != nil {
		return
	}
	h.subscribes[ch] = domain.NewSubscribe(from, to, id)
	return
}

func (h *huobiWs) sendSubscribeMsg(ch string, id int64) error {
	return h.conn.Write(h.ctx, websocket.MessageText, []byte(
		fmt.Sprintf("{\"sub\": \"%s\", \"id\":\"%d\"}", ch, id)),
	)
}

func (h *huobiWs) ListSubscribes() domain.Subscribes {
	s := make(domain.Subscribes, len(h.subscribes))
	h.subMu.RLock()
	defer h.subMu.RUnlock()
	for k, v := range h.subscribes {
		s[k] = v
	}
	return s
}

func gzipDecompress(r io.Reader) ([]byte, error) {
	r, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

func convertHuobiWsDataToDomain(from, to string, d *huobiWsData) *domain.Data {
	if d == nil {
		return nil
	}
	b, _ := json.Marshal(&domain.Raw{
		FromSymbol:     from,
		ToSymbol:       to,
		Open24Hour:     d.Tick.Open,
		Volume24Hour:   d.Tick.Amount,
		Volume24HourTo: d.Tick.Vol,
		High24Hour:     d.Tick.High,
		Price:          d.Tick.Bid,
		LastUpdate:     d.Ts,
		Supply:         float64(d.Tick.Count),
	})
	return &domain.Data{
		FromSymbol:     from,
		ToSymbol:       to,
		Open24Hour:     d.Tick.Open,
		Volume24Hour:   d.Tick.Amount,
		Low24Hour:      d.Tick.Low,
		High24Hour:     d.Tick.High,
		Price:          d.Tick.Bid,
		Supply:         float64(d.Tick.Count),
		LastUpdate:     d.Ts,
		DisplayDataRaw: string(b),
	}
}
