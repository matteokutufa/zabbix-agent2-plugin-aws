// File: factory/factory.go
package factory

import (
    "github.com/matteokutufa/zabbix-agent2-plugin-aws/aws"
    "github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// NewAWSClient crea un nuovo client AWS
func NewAWSClient(account models.AWSAccount) (models.AWSClientInterface, error) {
    return aws.NewClient(account)
}

// NewRDSDiscoverer crea un nuovo discoverer RDS
func NewRDSDiscoverer(client models.AWSClientInterface) models.RDSDiscovererInterface {
    return aws.NewRDSDiscoverer(client)
}

// NewMetricsCollector crea un nuovo raccoglitore di metriche
func NewMetricsCollector(client models.AWSClientInterface) models.MetricsCollectorInterface {
    return aws.NewMetricsCollector(client)
}

// LoadAccounts carica le configurazioni degli account
func LoadAccounts(path string) (map[string]models.AWSAccount, error) {
    return aws.LoadAccounts(path)
}

// LoadMetricsConfig carica la configurazione delle metriche
func LoadMetricsConfig(path string) (*models.MetricsConfig, error) {
    return aws.LoadMetricsConfig(path)
}