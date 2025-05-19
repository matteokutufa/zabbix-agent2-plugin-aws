// File: factory/factory.go
package factory

import (
    "github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// ProviderRegistry è un registro di provider di implementazioni
type ProviderRegistry struct {
    clientProvider          func(account models.AWSAccount) (models.AWSClientInterface, error)
    discovererProvider      func(client models.AWSClientInterface) models.RDSDiscovererInterface
    collectorProvider       func(client models.AWSClientInterface) models.MetricsCollectorInterface
    accountsLoaderProvider  func(path string) (map[string]models.AWSAccount, error)
    metricsConfigProvider   func(path string) (*models.MetricsConfig, error)
}

// DefaultProviderRegistry contiene i provider di default
var DefaultProviderRegistry ProviderRegistry

// RegisterClientProvider registra un provider per AWSClientInterface
func RegisterClientProvider(provider func(account models.AWSAccount) (models.AWSClientInterface, error)) {
    DefaultProviderRegistry.clientProvider = provider
}

// RegisterRDSDiscovererProvider registra un provider per RDSDiscovererInterface
func RegisterRDSDiscovererProvider(provider func(client models.AWSClientInterface) models.RDSDiscovererInterface) {
    DefaultProviderRegistry.discovererProvider = provider
}

// RegisterMetricsCollectorProvider registra un provider per MetricsCollectorInterface
func RegisterMetricsCollectorProvider(provider func(client models.AWSClientInterface) models.MetricsCollectorInterface) {
    DefaultProviderRegistry.collectorProvider = provider
}

// RegisterAccountsLoaderProvider registra un provider per il caricamento degli account
func RegisterAccountsLoaderProvider(provider func(path string) (map[string]models.AWSAccount, error)) {
    DefaultProviderRegistry.accountsLoaderProvider = provider
}

// RegisterMetricsConfigProvider registra un provider per il caricamento della configurazione delle metriche
func RegisterMetricsConfigProvider(provider func(path string) (*models.MetricsConfig, error)) {
    DefaultProviderRegistry.metricsConfigProvider = provider
}

// NewAWSClient crea un nuovo client AWS
func NewAWSClient(account models.AWSAccount) (models.AWSClientInterface, error) {
    if DefaultProviderRegistry.clientProvider == nil {
        panic("No client provider registered")
    }
    return DefaultProviderRegistry.clientProvider(account)
}

// NewRDSDiscoverer crea un nuovo discoverer RDS
func NewRDSDiscoverer(client models.AWSClientInterface) models.RDSDiscovererInterface {
    if DefaultProviderRegistry.discovererProvider == nil {
        panic("No RDS discoverer provider registered")
    }
    return DefaultProviderRegistry.discovererProvider(client)
}

// NewMetricsCollector crea un nuovo raccoglitore di metriche
func NewMetricsCollector(client models.AWSClientInterface) models.MetricsCollectorInterface {
    if DefaultProviderRegistry.collectorProvider == nil {
        panic("No metrics collector provider registered")
    }
    return DefaultProviderRegistry.collectorProvider(client)
}

// LoadAccounts carica le configurazioni degli account
func LoadAccounts(path string) (map[string]models.AWSAccount, error) {
    if DefaultProviderRegistry.accountsLoaderProvider == nil {
        panic("No accounts loader provider registered")
    }
    return DefaultProviderRegistry.accountsLoaderProvider(path)
}

// LoadMetricsConfig carica la configurazione delle metriche
func LoadMetricsConfig(path string) (*models.MetricsConfig, error) {
    if DefaultProviderRegistry.metricsConfigProvider == nil {
        panic("No metrics config provider registered")
    }
    return DefaultProviderRegistry.metricsConfigProvider(path)
}