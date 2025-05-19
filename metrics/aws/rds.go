// File: metrics/aws/rds.go
package aws

import (
	"encoding/json"
	"fmt"
	"time"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/zbxerr"

	"github.com/matteokutufa/zabbix-agent2-plugin-aws/factory"
	"github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// RDSDiscovery esegue il discovery delle istanze RDS
func RDSDiscovery(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
	if err := validateParams(params, 1); err != nil {
		return nil, err
	}

	accountID := params[0]

	// Ottieni il client AWS
	client, err := getClient(accountID)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	// Crea un discoverer RDS
	discoverer := factory.NewRDSDiscoverer(client)

	// Esegui il discovery
	result, err := discoverer.DiscoverInstances()
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	return result, nil
}

// RDSGet ottiene una metrica RDS
func RDSGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
	if err := validateParams(params, 3); err != nil {
		return nil, err
	}

	accountID := params[0]
	instanceID := params[1]
	metricName := params[2]

	// Opzionalmente, il parametro 3 può essere la statistica (default: Average)
	statistic := "Average"
	if len(params) > 3 && params[3] != "" {
		statistic = params[3]
	}

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

	// Trova la configurazione della metrica
	var metricConfig *models.MetricConfig
	for _, m := range metricsConfig.Services["rds"].Metrics {
		if m.Name == metricName {
			metricConfig = &m
			break
		}
	}

	if metricConfig == nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(fmt.Errorf("metric %s not found in configuration", metricName))
	}

	// Se è specificata una statistica nei parametri, usala invece di quella configurata
	if statistic != "Average" {
		metricConfig.Statistics = statistic
	}

	// Crea un collector per le metriche
	collector := factory.NewMetricsCollector(client)

	// Imposta l'orario di fine al momento attuale
	endTime := time.Now()

	// Imposta l'orario di inizio in base al periodo configurato
	startTime := endTime.Add(-time.Duration(metricConfig.Period) * time.Second)

	// Raccoglie la metrica
	value, err := collector.CollectRDSMetric(instanceID, metricName, metricConfig.Statistics, startTime, endTime)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	return value, nil
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