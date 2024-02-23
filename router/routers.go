package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/streamdp/ccd/clients"
	"github.com/streamdp/ccd/db"
	"github.com/streamdp/ccd/repos"
	"github.com/streamdp/ccd/router/handlers"
	v1 "github.com/streamdp/ccd/router/v1"
	"github.com/streamdp/ccd/router/v1/validators"
	"github.com/streamdp/ccd/router/v1/ws"
)

// InitRouter basic work on setting up the application, declare endpoints, register our custom validation functions
func InitRouter(
	e *gin.Engine,
	d db.Database,
	l *log.Logger,
	sr *repos.SymbolRepo,
	r clients.RestClient,
	w clients.WsClient,
	p clients.RestApiPuller,
) (err error) {
	// health checks
	e.GET("/healthz", SendOK)

	// serve web page
	e.LoadHTMLFiles("site/index.tmpl")
	e.Static("/css", "site/css")
	e.Static("/js", "site/js")
	e.GET("/", SendHTML)
	e.HEAD("/", SendOK)

	// serve api
	apiV1 := e.Group("/v1")
	{
		apiV1.GET("/collect/add", handlers.GinHandler(v1.AddWorker(p)))
		apiV1.GET("/collect/remove", handlers.GinHandler(v1.RemoveWorker(p)))
		apiV1.GET("/collect/update", handlers.GinHandler(v1.UpdateWorker(p)))
		apiV1.GET("/collect/status", handlers.GinHandler(v1.PullingStatus(p, w)))
		apiV1.GET("/symbols/add", handlers.GinHandler(v1.AddSymbol(sr)))
		apiV1.GET("/symbols/update", handlers.GinHandler(v1.UpdateSymbol(sr)))
		apiV1.GET("/symbols/remove", handlers.GinHandler(v1.RemoveSymbol(sr)))
		apiV1.GET("/price", handlers.GinHandler(v1.Price(r, d)))
		apiV1.GET("/ws", ws.HandleWs(r, l, d))

		apiV1.POST("/collect", handlers.GinHandler(v1.AddWorker(p)))
		apiV1.PUT("/collect", handlers.GinHandler(v1.UpdateWorker(p)))
		apiV1.DELETE("/collect", handlers.GinHandler(v1.RemoveWorker(p)))
		apiV1.POST("/symbols", handlers.GinHandler(v1.AddSymbol(sr)))
		apiV1.PUT("/symbols", handlers.GinHandler(v1.UpdateSymbol(sr)))
		apiV1.DELETE("/symbols", handlers.GinHandler(v1.RemoveSymbol(sr)))
		apiV1.POST("/price", handlers.GinHandler(v1.Price(r, d)))
		if w != nil {
			apiV1.POST("/ws/subscribe", handlers.GinHandler(v1.Subscribe(w)))
			apiV1.GET("/ws/subscribe", handlers.GinHandler(v1.Subscribe(w)))
			apiV1.POST("/ws/unsubscribe", handlers.GinHandler(v1.Unsubscribe(w)))
			apiV1.GET("/ws/unsubscribe", handlers.GinHandler(v1.Unsubscribe(w)))
		}
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err = v.RegisterValidation("symbols", validators.Symbols(sr)); err != nil {
			return err
		}
	}

	return nil
}
