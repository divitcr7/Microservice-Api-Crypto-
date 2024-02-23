package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamdp/ccd/clients"
	"github.com/streamdp/ccd/domain"
	"github.com/streamdp/ccd/router/handlers"
)

// CollectQuery structure for easily json serialization/validation/binding GET and POST query data
type CollectQuery struct {
	From     string `json:"fsym" form:"fsym" binding:"required,symbols"`
	To       string `json:"tsym" form:"tsym" binding:"required,symbols"`
	Interval int64  `json:"interval" form:"interval,default=60"`
}

// AddWorker that will collect data for the selected currency pair to the management service
func AddWorker(p clients.RestApiPuller) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		var t *clients.Task
		q := CollectQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		if t = p.Task(q.From, q.To); t != nil {
			r.UpdateAllFields(http.StatusOK, "Data for this pair is already being collected", t)
			return
		}
		t = p.AddTask(q.From, q.To, q.Interval)
		r.UpdateAllFields(http.StatusCreated, "Data collection started", t)
		return
	}
}

// RemoveWorker from the management service and stop collecting data for the selected currencies pair
func RemoveWorker(p clients.RestApiPuller) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		q := CollectQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		if p.Task(q.From, q.To) == nil {
			r.UpdateAllFields(http.StatusOK, "No data is collected for this pair", nil)
			return
		}
		p.RemoveTask(q.From, q.To)
		r.UpdateAllFields(http.StatusOK, "Task stopped successfully", nil)
		return
	}
}

// PullingStatus return information about running pull tasks
func PullingStatus(p clients.RestApiPuller, w clients.WsClient) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		r.UpdateAllFields(http.StatusOK, "Information about running tasks", nil)
		var (
			tasks      clients.Tasks
			subscribes domain.Subscribes
		)
		if p != nil {
			tasks = p.ListTasks()
		}
		if w != nil {
			subscribes = w.ListSubscribes()
		}
		if len(tasks) == 0 && len(subscribes) == 0 {
			return
		}
		list := map[string]map[string]interface{}{}
		for _, v := range tasks {
			if list[v.From] != nil {
				list[v.From][v.To] = v
				continue
			}
			list[v.From] = make(map[string]interface{})
			list[v.From][v.To] = v
		}
		for _, v := range subscribes {
			if list[v.From] != nil {
				list[v.From][v.To] = v
				continue
			}
			list[v.From] = make(map[string]interface{})
			list[v.From][v.To] = v
		}
		r.UpdateDataField(list)
		return
	}
}

// UpdateWorker update pulling data interval for the selected worker by the currencies pair
func UpdateWorker(p clients.RestApiPuller) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		var t *clients.Task
		q := CollectQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		if t = p.Task(q.From, q.To); t == nil {
			r.UpdateAllFields(http.StatusOK, "No data is collected for this pair", t)
			return
		}
		p.UpdateTask(t, q.Interval)
		r.UpdateAllFields(http.StatusOK, "Task updated successfully", t)
		return
	}
}

func Subscribe(w clients.WsClient) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		q := CollectQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		if err = w.Subscribe(q.From, q.To); err != nil {
			r.UpdateAllFields(http.StatusOK, "subscribe error:", err)
			return
		}
		r.UpdateAllFields(http.StatusCreated, "Subscribed successfully, data collection started", []string{q.From, q.To})
		return
	}
}

func Unsubscribe(w clients.WsClient) handlers.HandlerFuncResError {
	return func(c *gin.Context) (r handlers.Result, err error) {
		q := CollectQuery{}
		if err = c.Bind(&q); err != nil {
			return
		}
		if err = w.Unsubscribe(q.From, q.To); err != nil {
			r.UpdateAllFields(http.StatusOK, "Unsubscribe error:", err)
			return
		}
		r.UpdateAllFields(http.StatusOK, "Unsubscribed successfully, data collection stopped ", []string{q.From, q.To})
		return
	}
}
