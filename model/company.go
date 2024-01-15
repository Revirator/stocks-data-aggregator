package model

type Company struct {
	CIK        string                     `json:"cik"`
	Ticker     string                     `json:"ticker"`
	Name       string                     `json:"name"`
	Exchange   string                     `json:"exchange"`
	Financials map[string]FinancialMetric `json:"financials"`
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
