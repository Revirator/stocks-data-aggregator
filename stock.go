package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
)

type Stock struct {
	CIK        string      `json:"cik"`
	Ticker     string      `json:"ticker"`
	Name       string      `json:"name"`
	Exchange   string      `json:"exchange"`
	Financials *Financials `json:"financials"`
}

// TODO: change if needed
type Financials struct {
	Cash                                  FinancialsMetric
	CashAndCashEquivalentsAtCarryingValue FinancialsMetric
	CommonStockSharesOutstanding          FinancialsMetric
	CostsAndExpenses                      FinancialsMetric
	EarningsPerShareDiluted               FinancialsMetric
	LongTermDebt                          FinancialsMetric
	NetIncomeLoss                         FinancialsMetric
	PaymentsOfDividends                   FinancialsMetric
	PaymentsOfDividendsCommonStock        FinancialsMetric
	Revenues                              FinancialsMetric
	ShortTermInvestments                  FinancialsMetric
}

type FinancialsMetric struct {
	Description string            `json:"description"`
	Values      []FinancialsEntry `json:"values"`
}

type FinancialsEntry struct {
	Start string  `json:"start"`
	End   string  `json:"end"`
	Val   float64 `json:"val"`
	Form  string  `json:"form"`
	Frame string  `json:"frame"`
}

type DatabaseOperationOutcome string

const (
	SUCCESS             DatabaseOperationOutcome = "SUCCESS"
	DATABASE_ERROR      DatabaseOperationOutcome = "DATABASE_ERROR"
	STOCK_MISSING_ERROR DatabaseOperationOutcome = "STOCK_MISSING_ERROR"
)

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

func (database *Database) GetStockByTicker(ticker string) (*Stock, DatabaseOperationOutcome) {
	query := "SELECT * FROM STOCKS WHERE TICKER = $1"
	rows, err := database.db.Query(query, ticker)

	if err != nil {
		log.Printf(STOCK_RETRIEVAL_ERROR_LOG, ticker, err)
		return nil, DATABASE_ERROR
	}

	if !rows.Next() {
		log.Printf("Stock with ticker '%s' is missing in the database", ticker)
		return nil, STOCK_MISSING_ERROR
	}

	var financialData sql.NullString
	stock := &Stock{}
	err = rows.Scan(
		&stock.Ticker,
		&stock.CIK,
		&stock.Name,
		&stock.Exchange,
		&financialData,
	)
	if err != nil {
		log.Printf(STOCK_RETRIEVAL_ERROR_LOG, ticker, err)
		return nil, DATABASE_ERROR
	}

	if financialData.Valid {
		err = json.Unmarshal([]byte(financialData.String), &stock.Financials)
		if err != nil {
			log.Printf(STOCK_RETRIEVAL_ERROR_LOG, ticker, err)
			return nil, DATABASE_ERROR
		}
	}

	return stock, SUCCESS
}

func (database *Database) UpdateStockFinancialsByTicker(ticker string, financials *Financials) {
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
