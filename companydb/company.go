package companydb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"github.com/revirator/cfd/model"
)

type CompanyDatabase struct {
	DB *sql.DB
}

func (database CompanyDatabase) GetCompanyByTicker(ticker string) (*model.Company, error) {
	query := "SELECT * FROM COMPANIES WHERE TICKER = $1"
	rows, err := database.DB.Query(query, ticker)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("company missing")
	}

	var financialData sql.NullString
	company := &model.Company{}
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

func (database CompanyDatabase) UpdateCompanyFinancialsByTicker(
	ticker string, financials map[string]model.FinancialMetric,
) {
	financialData, err := json.Marshal(financials)
	if err != nil {
		log.Println(err)
		return
	}

	query := "UPDATE COMPANIES SET FINANCIALS = $1 WHERE TICKER = $2"
	_, err = database.DB.Exec(query, financialData, ticker)
	if err != nil {
		log.Println(err)
	}
}
