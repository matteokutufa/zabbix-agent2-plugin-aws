// File: aws/config.go
package aws

import (
    "io/ioutil"

    "gopkg.in/ini.v1"
    "gopkg.in/yaml.v2"

    "github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// LoadAccounts carica le configurazioni degli account dal file specificato
func LoadAccounts(filePath string) (map[string]models.AWSAccount, error) {
    // Implementazione di base
    result := make(map[string]models.AWSAccount)

    // Carica il file INI
    cfg, err := ini.Load(filePath)
    if err != nil {
        return nil, err
    }

    // Itera su tutte le sezioni (ogni sezione rappresenta un account)
    for _, section := range cfg.Sections() {
        // Salta la sezione DEFAULT
        if section.Name() == "DEFAULT" {
            continue
        }

        // Crea un nuovo account
        account := models.AWSAccount{
            Name:            section.Key("name").String(),
            AccessKey:       section.Key("access_key").String(),
            SecretAccessKey: section.Key("secret_access_key").String(),
            Region:          section.Key("region").String(),
            RoleARN:         section.Key("role_arn").String(),
        }

        // Aggiungi l'account al risultato
        result[section.Name()] = account
    }

    return result, nil
}

// LoadMetricsConfig carica la configurazione delle metriche dal file specificato
func LoadMetricsConfig(filePath string) (*models.MetricsConfig, error) {
    // Implementazione di base
    config := &models.MetricsConfig{
        Services: make(map[string]models.ServiceConfig),
    }

    // Leggi il file YAML
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    // Struttura temporanea per il parsing
    type tempMetric struct {
        Name       string `yaml:"name"`
        Statistics string `yaml:"statistics"`
        Period     int64  `yaml:"period"`
        Dimension  string `yaml:"dimension,omitempty"`
    }

    type tempDiscovery struct {
        Interval int64 `yaml:"interval"`
    }

    type tempService struct {
        Metrics   []tempMetric  `yaml:"metrics"`
        Discovery tempDiscovery `yaml:"discovery"`
    }

    type tempConfig struct {
        Services map[string]tempService `yaml:"services"`
    }

    // Parse il file YAML
    var tempCfg tempConfig
    if err := yaml.Unmarshal(data, &tempCfg); err != nil {
        return nil, err
    }

    // Converti la struttura temporanea nella struttura finale
    for serviceName, tempService := range tempCfg.Services {
        serviceConfig := models.ServiceConfig{
            Metrics: make([]models.MetricConfig, len(tempService.Metrics)),
        }

        for i, tempMetric := range tempService.Metrics {
            serviceConfig.Metrics[i] = models.MetricConfig{
                Name:       tempMetric.Name,
                Statistics: tempMetric.Statistics,
                Period:     tempMetric.Period,
                Dimension:  tempMetric.Dimension,
            }
        }

        serviceConfig.Discovery.Interval = tempService.Discovery.Interval

        config.Services[serviceName] = serviceConfig
    }

    return config, nil
}