package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamdp/ccd/repos"
	"github.com/streamdp/ccd/router/handlers"
)

// SymbolQuery structure for easily json serialization/validation/binding GET and POST query data
type SymbolQuery struct {
	Symbol  string `json:"symbol" form:"symbol" binding:"required"`
	Unicode string `json:"unicode" form:"unicode" binding:"required"`
}

// AddSymbol to the symbols table
func AddSymbol(sr *repos.SymbolRepo) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		q := SymbolQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		if err = sr.Add(q.Symbol, q.Unicode); err != nil {
			return
		}
		r.UpdateAllFields(http.StatusOK, fmt.Sprintf("symbol %s successfully added to the db", q.Symbol), nil)
		return
	}
}

// UpdateSymbol in the symbols table
func UpdateSymbol(sr *repos.SymbolRepo) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		q := SymbolQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		if err = sr.Update(q.Symbol, q.Unicode); err != nil {
			return
		}
		r.UpdateAllFields(http.StatusOK, fmt.Sprintf("symbol %s successfully updated", q.Symbol), nil)
		return
	}
}

// RemoveSymbol from the symbols table
func RemoveSymbol(sr *repos.SymbolRepo) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		q := SymbolQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		if err = sr.Remove(q.Symbol); err != nil {
			return
		}
		r.UpdateAllFields(http.StatusOK, fmt.Sprintf("symbol %s successfully removed", q.Symbol), nil)
		return
	}
}
