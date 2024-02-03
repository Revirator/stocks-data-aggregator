package main

import (
	"log"

	"github.com/revirator/cfd/external"
	"github.com/revirator/cfd/model"
)

func MapFinancialFactsToFinancialMetrics(facts *external.FinancialFacts) map[string]*model.FinancialMetric {
	return map[string]*model.FinancialMetric{
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

func mapMetricToFinancialMetric(fact *external.Metric) *model.FinancialMetric {
	if fact == nil {
		return nil
	}
	metric := model.FinancialMetric{}
	metric.Label = fact.Label
	metric.Description = fact.Description
	mapUnitsToFinancialEntries(fact.Units, &metric)
	return &metric
}

func mapUnitsToFinancialEntries(units external.Units, metric *model.FinancialMetric) {
	var entries []external.FinancialDataEntry
	if len(units.PrimaryEntries) > 0 {
		entries = units.PrimaryEntries
	} else if len(units.SecondaryEntries) > 0 {
		entries = units.SecondaryEntries
	} else if len(units.TertiaryEntries) > 0 {
		entries = units.TertiaryEntries
	} else {
		return
	}

	// TODO: handle 10-Q/A forms differently?
	// TODO: handle 8-K and other forms?
	for _, entry := range entries {
		if entry.IsQuarterlyReport() {
			metric.Quarterly = append(metric.Quarterly, mapFinancialDataEntryToFinancialEntry(entry))
		}
		if entry.IsAnnualReport() {
			metric.Annually = append(metric.Annually, mapFinancialDataEntryToFinancialEntry(entry))
		}
	}

	if metric.Label == "Revenues" {
		log.Println(metric.Annually)
	}
}

func mapFinancialDataEntryToFinancialEntry(entry external.FinancialDataEntry) model.FinancialEntry {
	return model.FinancialEntry{
		Value: entry.Value,
		Frame: entry.Frame,
		Form:  model.FinancialForm(entry.Form),
	}
}
