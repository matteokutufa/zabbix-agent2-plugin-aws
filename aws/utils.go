// File: aws/utils.go
package aws

// MetricsConfigFile è il percorso del file di configurazione delle metriche
var MetricsConfigFile string

// SetMetricsConfigFile imposta il percorso del file di configurazione delle metriche
func SetMetricsConfigFile(path string) {
    MetricsConfigFile = path
}

// AccountConfigFile è il percorso del file di configurazione degli account
var AccountConfigFile string

// SetAccountConfigFile imposta il percorso del file di configurazione degli account
func SetAccountConfigFile(path string) {
    AccountConfigFile = path
}