// File: aws/config.go
package aws

import (

    "io/ioutil"
    "path/filepath"

    "gopkg.in/ini.v1"
    "gopkg.in/yaml.v2"

	"github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// LoadAccounts carica le configurazioni degli account dal file specificato
func LoadAccounts(filePath string) (map[string]models.AWSAccount, error) {
	// Implementazione...
}

// LoadMetricsConfig carica la configurazione delle metriche dal file specificato
func LoadMetricsConfig(filePath string) (*models.MetricsConfig, error) {
	// Implementazione...
}
