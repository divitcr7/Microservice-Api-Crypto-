package cryptocompare

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/streamdp/ccd/clients"
	"github.com/streamdp/ccd/config"
	"github.com/streamdp/ccd/domain"
)

const (
	apiUrl = "https://min-api.cryptocompare.com"

	// Multiple Symbols Full Data - Get all the current trading info (price, vol, open, high, low etc) of any list of
	// cryptocurrencies in any other currency that you need. If the crypto does not trade directly into the toSymbol
	// requested, BTC will be used for conversion. This API also returns Display values for all the fields. If the
	// opposite pair trades we invert it (eg.: BTC-XMR)
	multipleSymbolsFullData = "/data/pricemultifull"
)

type cryptoCompareRest struct {
	apiKey string
	client *http.Client
}

// Init apiKey, apiUrl, wsURL variables with environment values and return CryptoCompareData structure
func Init() (rc clients.RestClient, err error) {
	var apiKey string
	if apiKey, err = getApiKey(); err != nil {
		return
	}
	return &cryptoCompareRest{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: time.Duration(config.HttpClientTimeout) * time.Millisecond,
		},
	}, nil
}

func getApiKey() (apiKey string, err error) {
	if apiKey = config.GetEnv("CCDC_APIKEY"); apiKey == "" {
		return "", errors.New("you should specify \"CCDC_APIKEY\" in you OS environment")
	}
	return
}

// Get filled CryptoCompareData structure for the selected pair currencies over http/https
func (cc *cryptoCompareRest) Get(fSym string, tSym string) (ds *domain.Data, err error) {
	var (
		u        *url.URL
		response *http.Response
		body     []byte
	)
	if u, err = cc.buildURL(fSym, tSym); err != nil {
		return nil, err
	}
	if response, err = cc.client.Get(u.String()); err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if body, err = io.ReadAll(response.Body); err != nil {
		return
	}
	if response.StatusCode != 200 {
		return
	}
	rawData := &cryptoCompareData{}
	if err = json.Unmarshal(body, rawData); err != nil {
		return
	}
	ds = convertToDomain(fSym, tSym, rawData)
	return
}

func convertToDomain(from, to string, d *cryptoCompareData) *domain.Data {
	r := d.Raw[from][to]
	b, _ := json.Marshal(&domain.Raw{
		FromSymbol:     from,
		ToSymbol:       to,
		Open24Hour:     r.Open24Hour,
		Volume24Hour:   r.Volume24Hour,
		Volume24HourTo: r.Volume24Hourto,
		High24Hour:     r.High24Hour,
		Price:          r.Price,
		LastUpdate:     r.LastUpdate,
		Supply:         r.Supply,
		MktCap:         r.MktCap,
	})
	return &domain.Data{
		FromSymbol:      from,
		ToSymbol:        to,
		Change24Hour:    r.Change24Hour,
		ChangePct24Hour: r.Changepct24Hour,
		Open24Hour:      r.Open24Hour,
		Volume24Hour:    r.Volume24Hour,
		Low24Hour:       r.Low24Hour,
		High24Hour:      r.High24Hour,
		Price:           r.Price,
		Supply:          r.Supply,
		MktCap:          r.MktCap,
		LastUpdate:      r.LastUpdate,
		DisplayDataRaw:  string(b),
	}
}

func (cc *cryptoCompareRest) buildURL(fSym string, tSym string) (u *url.URL, err error) {
	if u, err = url.Parse(apiUrl + multipleSymbolsFullData); err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("fsyms", fSym)
	query.Set("tsyms", tSym)
	query.Set("api_key", cc.apiKey)
	u.RawQuery = query.Encode()
	return u, nil
}
