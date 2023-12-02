package main

func MapFinancialFactsToFinancialMetrics(facts *FinancialFacts) map[string]FinancialMetric {
	return map[string]FinancialMetric{
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

func mapMetricToFinancialMetric(fact Metric) FinancialMetric {
	return FinancialMetric{
		Label:       fact.Label,
		Description: fact.Description,
		Values:      mapUnitsToFinancialEntries(fact.Units),
	}
}

func mapUnitsToFinancialEntries(units Units) []FinancialEntry {
	var entries []FinancialDataEntry
	if len(units.PrimaryEntries) > 0 {
		entries = units.PrimaryEntries
	} else if len(units.SecondaryEntries) > 0 {
		entries = units.SecondaryEntries
	} else if len(units.TertiaryEntries) > 0 {
		entries = units.TertiaryEntries
	} else {
		return []FinancialEntry{}
	}

	result := make([]FinancialEntry, len(entries))
	for _, entry := range entries {
		result = append(result, FinancialEntry(entry))
	}
	return result
}
