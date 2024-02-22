package huobi

type huobiRestData struct {
	Ch      string `json:"ch"`
	Status  string `json:"status"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
	Ts      int64  `json:"ts"`
	Tick    struct {
		Id      int64     `json:"id"`
		Version int64     `json:"version"`
		Open    float64   `json:"open"`
		Close   float64   `json:"close"`
		Low     float64   `json:"low"`
		High    float64   `json:"high"`
		Amount  float64   `json:"amount"`
		Vol     float64   `json:"vol"`
		Count   int       `json:"count"`
		Bid     []float64 `json:"bid"`
		Ask     []float64 `json:"ask"`
	} `json:"tick"`
}

type huobiWsData struct {
	Ch   string `json:"ch"`
	Ts   int64  `json:"ts"`
	Tick struct {
		Open      float64 `json:"open"`
		High      float64 `json:"high"`
		Low       float64 `json:"low"`
		Close     float64 `json:"close"`
		Amount    float64 `json:"amount"`
		Vol       float64 `json:"vol"`
		Count     int     `json:"count"`
		Bid       float64 `json:"bid"`
		BidSize   float64 `json:"bidSize"`
		Ask       float64 `json:"ask"`
		AskSize   float64 `json:"askSize"`
		LastPrice float64 `json:"lastPrice"`
		LastSize  float64 `json:"lastSize"`
	} `json:"tick"`
}
