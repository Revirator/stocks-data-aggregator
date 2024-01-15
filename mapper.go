package main

import (
	"github.com/revirator/cfd/external"
	"github.com/revirator/cfd/model"
)

func MapFinancialFactsToFinancialMetrics(facts *external.FinancialFacts) map[string]model.FinancialMetric {
	return map[string]model.FinancialMetric{
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

func mapMetricToFinancialMetric(fact external.Metric) model.FinancialMetric {
	return model.FinancialMetric{
		Label:       fact.Label,
		Description: fact.Description,
		Values:      mapUnitsToFinancialEntries(fact.Units),
	}
}

func mapUnitsToFinancialEntries(units external.Units) []model.FinancialEntry {
	var entries []external.FinancialDataEntry
	if len(units.PrimaryEntries) > 0 {
		entries = units.PrimaryEntries
	} else if len(units.SecondaryEntries) > 0 {
		entries = units.SecondaryEntries
	} else if len(units.TertiaryEntries) > 0 {
		entries = units.TertiaryEntries
	} else {
		return []model.FinancialEntry{}
	}

	result := make([]model.FinancialEntry, len(entries))
	for _, entry := range entries {
		result = append(result, model.FinancialEntry(entry))
	}
	return result
}
