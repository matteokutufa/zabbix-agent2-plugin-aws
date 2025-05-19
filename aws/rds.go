// rds.go
package aws

import (
	"fmt"
	"time"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/zbxerr"

	"github.com/yourname/zabbix-aws-plugins/aws"
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
	discoverer := aws.NewRDSDiscoverer(client)

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
	metricsConfig, err := aws.LoadMetricsConfig(MetricsFile())
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	// Trova la configurazione della metrica
	var metricConfig *aws.MetricConfig
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
	collector := aws.NewMetricsCollector(client)

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
