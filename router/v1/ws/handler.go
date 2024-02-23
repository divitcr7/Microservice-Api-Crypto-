package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streamdp/ccd/db"
	"github.com/streamdp/ccd/domain"
	"github.com/streamdp/ccd/router/handlers"
	"nhooyr.io/websocket"

	"github.com/streamdp/ccd/clients"
	v1 "github.com/streamdp/ccd/router/v1"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 512
)

type wsHandler struct {
	ctx         context.Context
	l           *log.Logger
	cancel      context.CancelFunc
	conn        *websocket.Conn
	messagePipe chan []byte

	rc clients.RestClient
	db db.Database
}

// HandleWs - handles websocket requests from the peer.
func HandleWs(r clients.RestClient, l *log.Logger, db db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(context.Background())
		conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			cancel()
			l.Println(err)
			return
		}
		h := &wsHandler{
			ctx:         ctx,
			l:           l,
			cancel:      cancel,
			conn:        conn,
			messagePipe: make(chan []byte, 256),
			rc:          r,
			db:          db,
		}
		h.conn.SetReadLimit(maxMessageSize)
		go h.handleMessagePipe()
		go h.handleClientRequests()
	}
}

func (w *wsHandler) handleClientRequests() {
	defer func() {
		w.cancel()
		close(w.messagePipe)
	}()
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			var (
				data  []byte
				err   error
				query = v1.PriceQuery{}
			)
			if _, data, err = w.conn.Read(w.ctx); err != nil {
				w.l.Println(err)
				if errors.As(err, &websocket.CloseError{}) {
					return
				}
				continue
			}
			if err = json.Unmarshal(data, &query); err != nil {
				w.returnAnErrorToTheClient(errors.New(
					"invalid request: the request should look like {\"fsym\":\"CRYPTO\",\"tsym\":\"COMMON\"}",
				))
				continue
			}
			if data, err = w.getLastPrice(&query); err != nil {
				w.l.Println(err)
				continue
			}
			w.messagePipe <- data
		}
	}
}

func (w *wsHandler) getLastPrice(q *v1.PriceQuery) (result []byte, err error) {
	var data *domain.Data
	if data, err = v1.LastPrice(w.rc, w.db, q); err != nil {
		return
	}
	if result, err = json.Marshal(&data); err != nil {
		return
	}
	return
}

func (w *wsHandler) handleMessagePipe() {
	defer w.cancel()
	for message := range w.messagePipe {
		ctx, cancel := context.WithTimeout(w.ctx, writeWait)
		if err := w.conn.Write(ctx, websocket.MessageText, message); err != nil {
			w.l.Println(err)
			cancel()
			return
		}
		cancel()
	}
	if err := w.conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		w.l.Println(err)
		return
	}
}

func (w *wsHandler) returnAnErrorToTheClient(err error) {
	var binaryString []byte
	r := handlers.Result{}
	r.UpdateAllFields(http.StatusBadRequest, err.Error(), nil)
	if binaryString, err = json.Marshal(&r); err != nil {
		w.l.Println(err)
		return
	}
	w.messagePipe <- binaryString
}
