package companydb

import (
	"database/sql"
	"encoding/json"
	"errors"

	_ "github.com/lib/pq"
)

type Company struct {
	CIK               string                     `json:"cik"`
	Ticker            string                     `json:"ticker"`
	Name              string                     `json:"name"`
	Exchange          string                     `json:"exchange"`
	Financials        map[string]FinancialMetric `json:"financials"`
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

type CompanyDatabase struct {
	db *sql.DB
}

func NewCompanyDatabse(db *sql.DB) CompanyDatabase {
	return CompanyDatabase{db}
}

func (database CompanyDatabase) GetCompanyByTicker(ticker string) (*Company, error) {
	query := "SELECT * FROM COMPANIES WHERE TICKER = $1"
	rows, err := database.db.Query(query, ticker)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("stock missing")
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
		return nil, err
	}

	if financialData.Valid {
		err = json.Unmarshal([]byte(financialData.String), &company.Financials)
		if err != nil {
			return nil, err
		}
	}

	return company, nil
}

func (database CompanyDatabase) UpdateCompanyFinancialsByTicker(ticker string, financials map[string]FinancialMetric) error {
	financialData, err := json.Marshal(financials)
	if err != nil {
		return err
	}

	query := "UPDATE COMPANIES SET FINANCIALS = $1 WHERE TICKER = $2"
	_, err = database.db.Exec(query, financialData, ticker)
	return err
}
