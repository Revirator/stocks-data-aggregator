package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type EdgarEntry struct {
	Facts FinancialFacts `json:"facts"`
}

type FinancialFacts struct {
	Concept Concept `json:"us-gaap"`
}

type Concept struct {
	Cash                                  Metric // TODO: check why some companies are missing this
	CashAndCashEquivalentsAtCarryingValue Metric
	CommonStockSharesOutstanding          Metric // TODO: might not be up to date?
	CostsAndExpenses                      Metric
	EarningsPerShareDiluted               Metric
	LongTermDebt                          Metric
	NetIncomeLoss                         Metric
	PaymentsOfDividends                   Metric
	PaymentsOfDividendsCommonStock        Metric
	Revenues                              Metric
	ShortTermInvestments                  Metric
}

type Metric struct {
	Description string                    `json:"description"`
	Wrapper     FinancialDataEntryWrapper `json:"units"`
}

type FinancialDataEntryWrapper struct {
	PrimaryEntries   []FinancialDataEntry `json:"usd"`
	SecondaryEntries []FinancialDataEntry `json:"shares"`
	TertiaryEntries  []FinancialDataEntry `json:"usd/shares"`
}

type FinancialDataEntry struct {
	Start string  `json:"start"`
	End   string  `json:"end"`
	Val   float64 `json:"val"`
	Form  string  `json:"form"`
	Frame string  `json:"frame"`
}

const (
	EDGAR_HOST             = "www.sec.gov"
	EDGAR_COMPANY_DATA_URL = "https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json"
)

const ERROR_LOG = "Could not retrieve financial data from Edgar for stock with CIK '%s'. Root cause:\n%s"

var Client = http.Client{Timeout: time.Second * 10}

func GetFinancialsForCompanyGivenCIK(cik string) *Concept {
	request, err := http.NewRequest("GET", fmt.Sprintf(EDGAR_COMPANY_DATA_URL, cik), nil)
	if err != nil {
		log.Printf(ERROR_LOG, cik, err)
		return nil
	}
	request.Header.Add("User-Agent", os.Getenv("EMAIl"))
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Host", EDGAR_HOST)

	response, err := Client.Do(request)
	if err != nil {
		log.Printf(ERROR_LOG, cik, err)
		return nil
	}
	defer response.Body.Close()

	var reader io.Reader
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			log.Printf(ERROR_LOG, cik, err)
			return nil
		}
		defer reader.(*gzip.Reader).Close()
	default:
		reader = response.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		log.Printf(ERROR_LOG, cik, err)
		return nil
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Could not retrieve financial data from Edgar for stock with CIK '%s'", cik)
		log.Printf("Got a response with statuc code %d and body:\n%s", response.StatusCode, string(body))
		return nil
	}

	data := EdgarEntry{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf(ERROR_LOG, cik, err)
		return nil
	}

	return &data.Facts.Concept
}
