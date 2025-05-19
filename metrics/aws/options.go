// options.go
package aws

import (
	"time"
)

// Options rappresenta le opzioni di configurazione del plugin
type Options struct {
	// Percorso del file di configurazione degli account AWS
	AccountFile string `conf:"optional,name=AccountFile"`

	// Percorso del file di configurazione delle metriche
	MetricsFile string `conf:"optional,name=MetricsFile"`

	// Timeout per le richieste API in secondi
	Timeout int `conf:"optional,range=1:30,name=Timeout"`

	// Keep-alive per le connessioni in secondi
	KeepAlive int `conf:"optional,range=60:900,name=KeepAlive"`

	// Numero massimo di sessioni concorrenti
	Sessions int `conf:"optional,range=1:100,name=Sessions"`
}

// DefaultOptions restituisce le opzioni di default
func DefaultOptions() Options {
	return Options{
		AccountFile: "/etc/zabbix/aws_accounts.ini",
		MetricsFile: "/etc/zabbix/metrics_config.yaml",
		Timeout:     30,
		KeepAlive:   300,
		Sessions:    15,
	}
}

// globalOptions contiene le opzioni di configurazione globali
var globalOptions Options

// Configure imposta le opzioni di configurazione globali
func Configure(options *Options) {
	if options != nil {
		globalOptions = *options
	} else {
		globalOptions = DefaultOptions()
	}
}

// Timeout restituisce il timeout configurato
func Timeout() time.Duration {
	return time.Duration(globalOptions.Timeout) * time.Second
}

// AccountFile restituisce il percorso del file degli account
func AccountFile() string {
	return globalOptions.AccountFile
}

// MetricsFile restituisce il percorso del file delle metriche
func MetricsFile() string {
	return globalOptions.MetricsFile
}

// Sessions restituisce il numero massimo di sessioni
func Sessions() int {
	return globalOptions.Sessions
}
