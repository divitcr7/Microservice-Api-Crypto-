package cryptocompare

type cryptoCompareData struct {
	Raw     map[string]map[string]*Response `json:"RAW"`
	Display map[string]map[string]*Display  `json:"DISPLAY"`
}

type cryptoCompareWsData struct {
	Type                string  `json:"TYPE"`
	FromSymbol          string  `json:"FROMSYMBOL"`
	ToSymbol            string  `json:"TOSYMBOL"`
	Price               float64 `json:"PRICE"`
	LastUpdate          int64   `json:"LASTUPDATE"`
	Volume24Hour        float64 `json:"VOLUME24HOUR"`
	Volume24HourTo      float64 `json:"VOLUME24HOURTO"`
	Open24Hour          float64 `json:"OPEN24HOUR"`
	High24Hour          float64 `json:"HIGH24HOUR"`
	Low24Hour           float64 `json:"LOW24HOUR"`
	CurrentSupply       float64 `json:"CURRENTSUPPLY"`
	CurrentSupplyMktCap float64 `json:"CURRENTSUPPLYMKTCAP"`
}

// Response structure for easily json serialization
type Response struct {
	Change24Hour    float64 `json:"CHANGE24HOUR"`
	Changepct24Hour float64 `json:"CHANGEPCT24HOUR"`
	Open24Hour      float64 `json:"OPEN24HOUR"`
	Volume24Hour    float64 `json:"VOLUME24HOUR"`
	Volume24Hourto  float64 `json:"VOLUME24HOURTO"`
	Low24Hour       float64 `json:"LOW24HOUR"`
	High24Hour      float64 `json:"HIGH24HOUR"`
	Price           float64 `json:"PRICE"`
	Supply          float64 `json:"SUPPLY"`
	MktCap          float64 `json:"MKTCAP"`
	LastUpdate      int64   `json:"LASTUPDATE"`
}

// Display structure for easily json serialization
type Display struct {
	Change24Hour    string `json:"CHANGE24HOUR"`
	Changepct24Hour string `json:"CHANGEPCT24HOUR"`
	Open24Hour      string `json:"OPEN24HOUR"`
	Volume24Hour    string `json:"VOLUME24HOUR"`
	Volume24Hourto  string `json:"VOLUME24HOURTO"`
	High24Hour      string `json:"HIGH24HOUR"`
	Price           string `json:"PRICE"`
	FromSymbol      string `json:"FROMSYMBOL"`
	ToSymbol        string `json:"TOSYMBOL"`
	LastUpdate      string `json:"LASTUPDATE"`
	Supply          string `json:"SUPPLY"`
	MktCap          string `json:"MKTCAP"`
}

// Data structure for easily json serialization
type Data struct {
	From    string
	To      string
	Raw     *Response `json:"RAW"`
	Display *Display  `json:"DISPLAY"`
}
