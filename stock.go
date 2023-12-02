package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Stock struct {
	CIK         string                     `json:"cik"`
	Ticker      string                     `json:"ticker"`
	CompanyName string                     `json:"company_name"`
	Exchange    string                     `json:"exchange"`
	Financials  map[string]FinancialMetric `json:"financials"`
}

type FinancialMetric struct {
	Label       string           `json:"label"`
	Description string           `json:"description"`
	Values      []FinancialEntry `json:"values"`
}

type FinancialEntry struct {
	Start string  `json:"start"`
	End   string  `json:"end"`
	Val   float64 `json:"val"`
	Form  string  `json:"form"`
	Frame string  `json:"frame"`
}

const (
	STOCK_RETRIEVAL_ERROR_LOG         = "Could not retrieve stock with ticker '%s' from the database. Root cause:\n%s"
	STOCK_FINANCIALS_UPDATE_ERROR_LOG = "Could not update financials for stock with ticker '%s'. Root cause:\n%s"
)

type Database struct {
	db *sql.DB
}

func DatabaseInit(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (database *Database) GetStockByTicker(ticker string) (*Stock, *ServerError) {
	query := "SELECT * FROM STOCKS WHERE TICKER = $1"
	rows, err := database.db.Query(query, ticker)

	if err != nil {
		log.Printf(STOCK_RETRIEVAL_ERROR_LOG, ticker, err)
		return nil, InternalServerError()
	}

	if !rows.Next() {
		return nil, &ServerError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("Stock with ticker '%s' does not exist or is not listed on any of the US exchanges.", ticker),
		}
	}

	var financialData sql.NullString
	stock := &Stock{}
	err = rows.Scan(
		&stock.Ticker,
		&stock.CIK,
		&stock.CompanyName,
		&stock.Exchange,
		&financialData,
	)
	if err != nil {
		log.Printf(STOCK_RETRIEVAL_ERROR_LOG, ticker, err)
		return nil, InternalServerError()
	}

	if financialData.Valid {
		err = json.Unmarshal([]byte(financialData.String), &stock.Financials)
		if err != nil {
			log.Printf(STOCK_RETRIEVAL_ERROR_LOG, ticker, err)
			return nil, InternalServerError()
		}
	}

	return stock, nil
}

func (database *Database) UpdateStockFinancialsByTicker(ticker string, financials map[string]FinancialMetric) {
	financialData, err := json.Marshal(financials)
	if err != nil {
		log.Printf(STOCK_FINANCIALS_UPDATE_ERROR_LOG, ticker, err)
		return
	}

	query := "UPDATE STOCKS SET FINANCIALS = $1 WHERE TICKER = $2"
	_, err = database.db.Exec(query, financialData, ticker)
	if err != nil {
		log.Printf(STOCK_FINANCIALS_UPDATE_ERROR_LOG, ticker, err)
	}
}
