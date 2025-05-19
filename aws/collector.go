// Aggiornamento di collector.go per usare models
// File: aws/collector.go
package aws

import (
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/cloudwatch"

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
    // Periodo di tempo in secondi
    period := int64(60) // Un minuto è un buon default per le metriche CloudWatch

    // Ottiene la metrica mediante il client AWS
    return c.client.GetRDSMetric(instanceID, metricName, statistic, period, startTime, endTime)
}