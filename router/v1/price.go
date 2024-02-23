package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/streamdp/ccd/clients"
	"github.com/streamdp/ccd/db"
	"github.com/streamdp/ccd/domain"
	"github.com/streamdp/ccd/router/handlers"
)

// PriceQuery structure for easily json serialization/validation/binding GET and POST query data
type PriceQuery struct {
	From string `json:"fsym" form:"fsym" binding:"required,symbols"`
	To   string `json:"tsym" form:"tsym" binding:"required,symbols"`
}

// LastPrice return up-to-date data for the selected currencies pair
func LastPrice(r clients.RestClient, db db.Database, query *PriceQuery) (d *domain.Data, err error) {
	from, to := strings.ToUpper(query.From), strings.ToUpper(query.To)
	if d, err = r.Get(from, to); err != nil {
		return db.GetLast(from, to)
	}
	db.DataPipe() <- d
	return
}

// Price return up-to-date or most recent data for the selected currencies pair
func Price(rc clients.RestClient, db db.Database) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		q := PriceQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		p, err := LastPrice(rc, db, &q)
		if err != nil {
			return
		}
		r.UpdateAllFields(http.StatusOK, fmt.Sprintf("Most recent price, updated at %d", p.LastUpdate), p)
		return
	}
}
