// File: aws/collector.go (completo)
package aws

import (
    "time"

    "github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// MetricsCollector rappresenta un raccoglitore di metriche
type MetricsCollector struct {
    client models.AWSClientInterface
}

// NewMetricsCollector crea un nuovo raccoglitore di metriche
func NewMetricsCollector(client models.AWSClientInterface) *MetricsCollector {
    return &MetricsCollector{
        client: client,
    }
}

// CollectRDSMetric raccoglie una metrica RDS
func (c *MetricsCollector) CollectRDSMetric(instanceID, metricName, statistic string, startTime, endTime time.Time) (float64, error) {
    // Periodo di tempo in secondi (default: 60 secondi)
    period := int64(60)

    // Ottiene la metrica mediante il client AWS
    return c.client.GetRDSMetric(instanceID, metricName, statistic, period, startTime, endTime)
}