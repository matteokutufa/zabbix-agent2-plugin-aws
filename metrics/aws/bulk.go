// bulk.go
package aws

import (
	"encoding/json"
	"time"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/zbxerr"

	"github.com/matteokutufa/zabbix-agent2-plugin-aws/factory"
	"github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// MetricResult rappresenta il risultato di una metrica
type MetricResult struct {
	MetricName string  `json:"metric"`
	Value      float64 `json:"value"`
	Timestamp  int64   `json:"timestamp"`
	Statistic  string  `json:"statistic"`
}

// BulkResult rappresenta il risultato di una richiesta bulk
type BulkResult struct {
	ResourceID string         `json:"resource_id"`
	Metrics    []MetricResult `json:"metrics"`
}

// RDSBulkGet ottiene tutte le metriche RDS in un'unica chiamata
func RDSBulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
	if err := validateParams(params, 2); err != nil {
		return nil, err
	}

	accountID := params[0]
	instanceID := params[1]

	// Ottieni il client AWS
	client, err := getClient(accountID)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	// Carica la configurazione delle metriche
	metricsConfig, err := factory.LoadMetricsConfig(MetricsFile())
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	// Crea un collector per le metriche
	collector := factory.NewMetricsCollector(client)

	// Imposta l'orario di fine al momento attuale
	endTime := time.Now()

	// Crea un risultato bulk
	result := BulkResult{
		ResourceID: instanceID,
		Metrics:    make([]MetricResult, 0),
	}

	// Raccoglie tutte le metriche configurate per RDS
	for _, metricConfig := range metricsConfig.Services["rds"].Metrics {
		// Imposta l'orario di inizio in base al periodo configurato
		startTime := endTime.Add(-time.Duration(metricConfig.Period) * time.Second)

		// Raccoglie la metrica
		value, err := collector.CollectRDSMetric(instanceID, metricConfig.Name, metricConfig.Statistics, startTime, endTime)
		if err != nil {
			// Non fallire l'intera richiesta per una singola metrica fallita
			// ma aggiungi un messaggio di log
			continue
		}

		// Aggiungi il risultato
		result.Metrics = append(result.Metrics, MetricResult{
			MetricName: metricConfig.Name,
			Value:      value,
			Timestamp:  endTime.Unix(),
			Statistic:  metricConfig.Statistics,
		})
	}

	// Restituisci il risultato come JSON
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	return string(jsonResult), nil
}

// S3BulkGet ottiene tutte le metriche S3 in un'unica chiamata
//func S3BulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
//	// Implementazione simile per S3
//	// ...
//
//
// ELBBulkGet ottiene tutte le metriche ELB in un'unica chiamata
//func ELBBulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
//	// Implementazione simile per ELB
//	// ...
//
//
/// ECSBulkGet ottiene tutte le metriche ECS in un'unica chiamata
//func ECSBulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
//	// Implementazione simile per ECS
//	// ...
//
//
/// MSKBulkGet ottiene tutte le metriche MSK in un'unica chiamata
//func MSKBulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
//	// Implementazione simile per MSK
//	// ...
//
//