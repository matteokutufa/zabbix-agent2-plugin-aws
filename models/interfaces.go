package models

import (
    "time"
)

// AWSClientInterface definisce l'interfaccia per le operazioni del client AWS
type AWSClientInterface interface {
    GetRDSMetric(instanceID, metricName, statistic string, period int64, startTime, endTime time.Time) (float64, error)
    ListAvailableRDSMetrics(instanceID string) ([]CloudWatchMetric, error)
    // Aggiungi altri metodi necessari
}

// MetricsCollectorInterface definisce l'interfaccia per la raccolta delle metriche
type MetricsCollectorInterface interface {
    CollectRDSMetric(instanceID, metricName, statistic string, startTime, endTime time.Time) (float64, error)
    // Aggiungi altri metodi necessari
}

// RDSDiscovererInterface definisce l'interfaccia per il discovery delle istanze RDS
type RDSDiscovererInterface interface {
    DiscoverInstances() (interface{}, error)
}

// ConfigLoaderInterface definisce l'interfaccia per il caricamento delle configurazioni
type ConfigLoaderInterface interface {
    LoadAccounts(path string) (map[string]AWSAccount, error)
    LoadMetricsConfig(path string) (*MetricsConfig, error)
}