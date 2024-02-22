package huobi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/streamdp/ccd/clients"
	"github.com/streamdp/ccd/config"
	"github.com/streamdp/ccd/domain"
)

const (
	apiUrl = "https://api.huobi.pro"

	// Get Latest Aggregated Ticker https://huobiapi.github.io/docs/spot/v1/en/#get-latest-aggregated-ticker
	// This endpoint retrieves the latest ticker with some important 24h aggregated market data.
	// Request Parameters "symbol" (all supported trading symbol, e.g. btcusdt, bccbtc. Refer to /v1/common/symbols)
	latestAggregatedTicker = "/market/detail/merged"
)

type huobiRest struct {
	client *http.Client
}

func Init() (clients.RestClient, error) {
	return &huobiRest{
		client: &http.Client{
			Timeout: time.Duration(config.HttpClientTimeout) * time.Millisecond,
		},
	}, nil
}

func (h *huobiRest) Get(fSym string, tSym string) (ds *domain.Data, err error) {
	var (
		u        *url.URL
		response *http.Response
		body     []byte
	)
	if u, err = h.buildURL(fSym, tSym); err != nil {
		return nil, err
	}

	if response, err = h.client.Get(u.String()); err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if body, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, err
	}
	rawData := &huobiRestData{}
	if err = json.Unmarshal(body, rawData); err != nil {
		return nil, err
	}
	if rawData.Status == "error" {
		return nil, errors.New(rawData.ErrMsg)
	}
	return convertHuobiRestDataToDomain(fSym, tSym, rawData), nil
}

func convertHuobiRestDataToDomain(from, to string, d *huobiRestData) *domain.Data {
	if d == nil {
		return nil
	}
	var price float64
	if len(d.Tick.Bid) > 0 {
		price = d.Tick.Bid[0]
	}
	b, _ := json.Marshal(&domain.Raw{
		FromSymbol:     from,
		ToSymbol:       to,
		Open24Hour:     d.Tick.Open,
		Volume24Hour:   d.Tick.Amount,
		Volume24HourTo: d.Tick.Vol,
		High24Hour:     d.Tick.High,
		Price:          price,
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
		Price:          price,
		Supply:         float64(d.Tick.Count),
		LastUpdate:     d.Ts,
		DisplayDataRaw: string(b),
	}
}

func (h *huobiRest) buildURL(fSym string, tSym string) (u *url.URL, err error) {
	if u, err = url.Parse(apiUrl + latestAggregatedTicker); err != nil {
		return nil, err
	}
	if strings.ToLower(tSym) == "usd" {
		tSym = "usdt"
	}
	query := u.Query()
	query.Set("symbol", strings.ToLower(fSym+tSym))
	u.RawQuery = query.Encode()
	return u, nil
}
