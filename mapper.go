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
	result := []FinancialEntry{}
	if len(units.PrimaryEntries) > 0 {
		for _, entry := range units.PrimaryEntries {
			result = append(result, FinancialEntry(entry))
		}
		return result
	}

	if len(units.SecondaryEntries) > 0 {
		for _, entry := range units.SecondaryEntries {
			result = append(result, FinancialEntry(entry))
		}
		return result
	}

	if len(units.TertiaryEntries) > 0 {
		for _, entry := range units.TertiaryEntries {
			result = append(result, FinancialEntry(entry))
		}
		return result
	}

	return result
}
