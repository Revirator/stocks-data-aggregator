package clients

import (
	"encoding/json"
	"fmt"
	"log"
)

type YahooEntry struct {
	Chart YahooChart `json:"chart"`
}

type YahooChart struct {
	Result []YahooMetadata `json:"result"`
}

type YahooMetadata struct {
	Metadata CompanyMetadata `json:"meta"`
}

type CompanyMetadata struct {
	RegularMarketPrice float64 `json:"regularMarketPrice"`
	ChartPreviousClose float64 `json:"chartPreviousClose"`
}

const YAHOO_COMPANY_METADATA_URL = "https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d"

func GetCompanyMetadataGivenTicker(ticker string) *CompanyMetadata {
	request := PrepareRequest("GET", fmt.Sprintf(YAHOO_COMPANY_METADATA_URL, ticker))
	body := SendRequestAndGetBody(request)

	data := YahooEntry{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
		return nil
	}

	return &data.Chart.Result[0].Metadata
}
