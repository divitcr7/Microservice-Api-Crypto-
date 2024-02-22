package domain

type Data struct {
	Id              int64   `json:"id" db:"_id"`
	FromSymbol      string  `json:"from_sym" db:"fromSym"`
	ToSymbol        string  `json:"to_sym" db:"toSym"`
	Change24Hour    float64 `json:"change_24_hour" db:"change24hour"`
	ChangePct24Hour float64 `json:"change_pct_24_hour" db:"changepct24hour"`
	Open24Hour      float64 `json:"open_24_hour" db:"open24hour"`
	Volume24Hour    float64 `json:"volume_24_hour" db:"volume24hour"`
	Low24Hour       float64 `json:"low_24_hour" db:"low24hour"`
	High24Hour      float64 `json:"high_24_hour" db:"high24hour"`
	Price           float64 `json:"price" db:"price"`
	Supply          float64 `json:"supply" db:"supply"`
	MktCap          float64 `json:"mkt_cap" db:"mktcap"`
	LastUpdate      int64   `json:"last_update" db:"lastupdate"`
	DisplayDataRaw  string  `json:"display_data_raw" db:"displaydataraw"`
}

// Raw structure for easily json serialization
type Raw struct {
	FromSymbol      string  `json:"from_symbol"`
	ToSymbol        string  `json:"to_symbol"`
	Change24Hour    float64 `json:"change_24_hour"`
	ChangePct24Hour float64 `json:"changepct_24_hour"`
	Open24Hour      float64 `json:"open_24_hour"`
	Volume24Hour    float64 `json:"volume_24_hour"`
	Volume24HourTo  float64 `json:"volume_24_hour_to"`
	Low24Hour       float64 `json:"low_24_hour"`
	High24Hour      float64 `json:"high_24_hour"`
	Price           float64 `json:"price"`
	Supply          float64 `json:"supply"`
	MktCap          float64 `json:"mkt_cap"`
	LastUpdate      int64   `json:"last_update"`
}
