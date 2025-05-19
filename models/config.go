// Esempio di ristrutturazione del progetto per risolvere il ciclo di importazioni

// 1. Crea un pacchetto models per le strutture dati condivise
// File: models/config.go
package models

// MetricConfig rappresenta la configurazione di una metrica
type MetricConfig struct {
	Name       string
	Statistics string
	Period     int64
	Dimension  string
}

// ServiceConfig rappresenta la configurazione di un servizio
type ServiceConfig struct {
	Metrics   []MetricConfig
	Discovery struct {
		Interval int64
	}
}

// MetricsConfig rappresenta la configurazione completa delle metriche
type MetricsConfig struct {
	Services map[string]ServiceConfig
}

// AWSAccount rappresenta la configurazione di un account AWS
type AWSAccount struct {
	Name            string
	AccessKey       string
	SecretAccessKey string
	Region          string
	RoleARN         string
}
