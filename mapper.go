package main

func MapConceptToFinancials(concept *Concept) *Financials {
	return &Financials{
		Cash:                                  mapMetricToFinancialsMetric(&concept.Cash),
		CashAndCashEquivalentsAtCarryingValue: mapMetricToFinancialsMetric(&concept.CashAndCashEquivalentsAtCarryingValue),
		CommonStockSharesOutstanding:          mapMetricToFinancialsMetric(&concept.CommonStockSharesOutstanding),
		CostsAndExpenses:                      mapMetricToFinancialsMetric(&concept.CostsAndExpenses),
		EarningsPerShareDiluted:               mapMetricToFinancialsMetric(&concept.EarningsPerShareDiluted),
		LongTermDebt:                          mapMetricToFinancialsMetric(&concept.LongTermDebt),
		NetIncomeLoss:                         mapMetricToFinancialsMetric(&concept.NetIncomeLoss),
		PaymentsOfDividends:                   mapMetricToFinancialsMetric(&concept.PaymentsOfDividends),
		PaymentsOfDividendsCommonStock:        mapMetricToFinancialsMetric(&concept.PaymentsOfDividendsCommonStock),
		Revenues:                              mapMetricToFinancialsMetric(&concept.Revenues),
		ShortTermInvestments:                  mapMetricToFinancialsMetric(&concept.ShortTermInvestments),
	}
}

func mapMetricToFinancialsMetric(fact *Metric) FinancialsMetric {
	return FinancialsMetric{
		Description: fact.Description,
		Values:      mapFinancialDataEntryWrapperToFinancialsEntries(&fact.Wrapper),
	}
}

func mapFinancialDataEntryWrapperToFinancialsEntries(wrapper *FinancialDataEntryWrapper) []FinancialsEntry {
	result := []FinancialsEntry{}
	if len(wrapper.PrimaryEntries) > 0 {
		for _, v := range wrapper.PrimaryEntries {
			result = append(result, mapFinancialDataEntryToFinancialsEntry(&v))
		}
		return result
	}

	if len(wrapper.SecondaryEntries) > 0 {
		for _, v := range wrapper.SecondaryEntries {
			result = append(result, mapFinancialDataEntryToFinancialsEntry(&v))
		}
		return result
	}

	if len(wrapper.TertiaryEntries) > 0 {
		for _, v := range wrapper.TertiaryEntries {
			result = append(result, mapFinancialDataEntryToFinancialsEntry(&v))
		}
		return result
	}

	return result
}

func mapFinancialDataEntryToFinancialsEntry(data *FinancialDataEntry) FinancialsEntry {
	return FinancialsEntry{
		Start: data.Start,
		End:   data.End,
		Val:   data.Val,
		Form:  data.Form,
		Frame: data.Frame,
	}
}
