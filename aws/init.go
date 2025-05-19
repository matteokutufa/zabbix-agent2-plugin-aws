// File: aws/init.go
package aws

import (
    "github.com/matteokutufa/zabbix-agent2-plugin-aws/factory"
    "github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// Registra i provider di implementazione nel registry
func init() {
    // Registra il provider per il client AWS
    factory.RegisterClientProvider(func(account models.AWSAccount) (models.AWSClientInterface, error) {
        return NewClient(account)
    })

    // Registra il provider per il RDS discoverer
    factory.RegisterRDSDiscovererProvider(func(client models.AWSClientInterface) models.RDSDiscovererInterface {
        return NewRDSDiscoverer(client)
    })

    // Registra il provider per il metrics collector
    factory.RegisterMetricsCollectorProvider(func(client models.AWSClientInterface) models.MetricsCollectorInterface {
        return NewMetricsCollector(client)
    })

    // Registra il provider per il caricamento degli account
    factory.RegisterAccountsLoaderProvider(func(path string) (map[string]models.AWSAccount, error) {
        return LoadAccounts(path)
    })

    // Registra il provider per il caricamento della configurazione delle metriche
    factory.RegisterMetricsConfigProvider(func(path string) (*models.MetricsConfig, error) {
        return LoadMetricsConfig(path)
    })
}