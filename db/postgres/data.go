package postgres

import (
	"database/sql"
	"errors"

	"github.com/streamdp/ccd/domain"
)

// GetLast row with the most recent data for the selected currencies pair
func (d *Db) GetLast(from string, to string) (result *domain.Data, err error) {
	result = &domain.Data{
		FromSymbol: from,
		ToSymbol:   to,
	}
	query := `
		select 
		       _id,
		       change24hour,
		       changepct24hour,
		       open24hour,
		       volume24hour,
		       low24hour,
		       high24hour, 
		       price, 
		       supply,
		       mktcap, 
		       lastupdate,
		       displaydataraw
		from data 
		where fromSym=(select _id from symbols where symbol=$1)
		  and toSym=(select _id from symbols where symbol=$2)
		ORDER BY lastupdate DESC limit 1;
`
	err = d.QueryRow(query, from, to).Scan(
		&result.Id,
		&result.Change24Hour,
		&result.ChangePct24Hour,
		&result.Open24Hour,
		&result.Volume24Hour,
		&result.Low24Hour,
		&result.High24Hour,
		&result.Price,
		&result.Supply,
		&result.MktCap,
		&result.LastUpdate,
		&result.DisplayDataRaw,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Insert clients.Data from the clients.DataPipe to the Db
func (d *Db) Insert(data *domain.Data) (result sql.Result, err error) {
	if data == nil {
		return nil, errors.New("cant insert empty data")
	}
	if err != nil {
		return nil, err
	}
	query := `insert into data (
                  fromSym,
                  toSym,
                  change24hour,
                  changepct24hour,
                  open24hour,
                  volume24hour,
                  low24hour,
                  high24hour, 
                  price,
                  supply, 
                  mktcap, 
                  lastupdate,
                  displaydataraw
        ) 
		values (
		        (SELECT _id FROM symbols WHERE symbol=$1),
		        (SELECT _id FROM symbols WHERE symbol=$2),
		        $3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		)
`
	return d.Exec(
		query,
		&data.FromSymbol,
		&data.ToSymbol,
		&data.Change24Hour,
		&data.ChangePct24Hour,
		&data.Open24Hour,
		&data.Volume24Hour,
		&data.Low24Hour,
		&data.High24Hour,
		&data.Price,
		&data.Supply,
		&data.MktCap,
		&data.LastUpdate,
		&data.DisplayDataRaw,
	)
}
