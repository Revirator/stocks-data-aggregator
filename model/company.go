package model

type Company struct {
	CIK        string                      `json:"cik"`
	Ticker     string                      `json:"ticker"`
	Name       string                      `json:"name"`
	Exchange   *string                     `json:"exchange"`
	Financials map[string]*FinancialMetric `json:"financials"`
}

type FinancialMetric struct {
	Label       string           `json:"label"`
	Description string           `json:"description"`
	Annually    []FinancialEntry `json:"annually"`
	Quarterly   []FinancialEntry `json:"quarterly"`
}

type FinancialEntry struct {
	Value float64       `json:"value"`
	Frame string        `json:"frame"`
	Form  FinancialForm `json:"form"`
}

type FinancialForm string

const (
	Q10 FinancialForm = "10-Q"
	K10 FinancialForm = "10-K"
)
