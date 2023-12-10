package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Company struct {
	CIK               string                     `json:"cik"`
	Ticker            string                     `json:"ticker"`
	Name              string                     `json:"name"`
	Exchange          string                     `json:"exchange"`
	Financials        map[string]FinancialMetric `json:"financials"`
	StockPrice        float64                    `json:"stock_price"`
	DayMovePercentage float64                    `json:"day_move_percentage"`
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
	COMPANY_RETRIEVAL_ERROR_LOG         = "Could not retrieve company with ticker '%s' from the database. Root cause:\n%s"
	COMPANY_FINANCIALS_UPDATE_ERROR_LOG = "Could not update financials for company with ticker '%s'. Root cause:\n%s"
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

func (database *Database) GetCompanyByTicker(ticker string) (*Company, *ServerError) {
	query := "SELECT * FROM COMPANIES WHERE TICKER = $1"
	rows, err := database.db.Query(query, ticker)

	if err != nil {
		log.Printf(COMPANY_RETRIEVAL_ERROR_LOG, ticker, err)
		return nil, InternalServerError()
	}

	if !rows.Next() {
		return nil, &ServerError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("Company with ticker '%s' does not exist or is not listed on any of the US exchanges.", ticker),
		}
	}

	var financialData sql.NullString
	company := &Company{}
	err = rows.Scan(
		&company.Ticker,
		&company.CIK,
		&company.Name,
		&company.Exchange,
		&financialData,
	)
	if err != nil {
		log.Printf(COMPANY_RETRIEVAL_ERROR_LOG, ticker, err)
		return nil, InternalServerError()
	}

	if financialData.Valid {
		err = json.Unmarshal([]byte(financialData.String), &company.Financials)
		if err != nil {
			log.Printf(COMPANY_RETRIEVAL_ERROR_LOG, ticker, err)
			return nil, InternalServerError()
		}
	}

	return company, nil
}

func (database *Database) UpdateCompanyFinancialsByTicker(ticker string, financials map[string]FinancialMetric) {
	financialData, err := json.Marshal(financials)
	if err != nil {
		log.Printf(COMPANY_FINANCIALS_UPDATE_ERROR_LOG, ticker, err)
		return
	}

	query := "UPDATE COMPANIES SET FINANCIALS = $1 WHERE TICKER = $2"
	_, err = database.db.Exec(query, financialData, ticker)
	if err != nil {
		log.Printf(COMPANY_FINANCIALS_UPDATE_ERROR_LOG, ticker, err)
	}
}
