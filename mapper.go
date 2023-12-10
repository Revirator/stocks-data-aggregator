package main

import (
	"github.com/revirator/cfd/clients"
	"github.com/revirator/cfd/companydb"
)

func MapFinancialFactsToFinancialMetrics(facts *clients.FinancialFacts) map[string]companydb.FinancialMetric {
	return map[string]companydb.FinancialMetric{
		"Cash":                                  mapMetricToFinancialMetric(facts.Principles.Cash),
		"CashAndCashEquivalentsAtCarryingValue": mapMetricToFinancialMetric(facts.Principles.CashAndCashEquivalentsAtCarryingValue),
		"CommonStockSharesOutstanding":          mapMetricToFinancialMetric(facts.Principles.CommonStockSharesOutstanding),
		"CostsAndExpenses":                      mapMetricToFinancialMetric(facts.Principles.CostsAndExpenses),
		"EarningsPerShareDiluted":               mapMetricToFinancialMetric(facts.Principles.EarningsPerShareDiluted),
		"EntityCommonStockSharesOutstanding":    mapMetricToFinancialMetric(facts.Entity.EntityCommonStockSharesOutstanding),
		"LongTermDebt":                          mapMetricToFinancialMetric(facts.Principles.LongTermDebt),
		"NetIncomeLoss":                         mapMetricToFinancialMetric(facts.Principles.NetIncomeLoss),
		"PaymentsOfDividends":                   mapMetricToFinancialMetric(facts.Principles.PaymentsOfDividends),
		"PaymentsOfDividendsCommonStock":        mapMetricToFinancialMetric(facts.Principles.PaymentsOfDividendsCommonStock),
		"Revenues":                              mapMetricToFinancialMetric(facts.Principles.Revenues),
		"ShortTermInvestments":                  mapMetricToFinancialMetric(facts.Principles.ShortTermInvestments),
	}
}

func mapMetricToFinancialMetric(fact clients.Metric) companydb.FinancialMetric {
	return companydb.FinancialMetric{
		Label:       fact.Label,
		Description: fact.Description,
		Values:      mapUnitsToFinancialEntries(fact.Units),
	}
}

func mapUnitsToFinancialEntries(units clients.Units) []companydb.FinancialEntry {
	var entries []clients.FinancialDataEntry
	if len(units.PrimaryEntries) > 0 {
		entries = units.PrimaryEntries
	} else if len(units.SecondaryEntries) > 0 {
		entries = units.SecondaryEntries
	} else if len(units.TertiaryEntries) > 0 {
		entries = units.TertiaryEntries
	} else {
		return []companydb.FinancialEntry{}
	}

	result := make([]companydb.FinancialEntry, len(entries))
	for _, entry := range entries {
		result = append(result, companydb.FinancialEntry(entry))
	}
	return result
}
