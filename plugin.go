// plugin.go asdasd
package main

import (
    "git.zabbix.com/ap/plugin-support/conf"
    "git.zabbix.com/ap/plugin-support/plugin"
    "git.zabbix.com/ap/plugin-support/plugin/comms"
    "git.zabbix.com/ap/plugin-support/zbxerr"

    // Importante: importa il pacchetto aws per l'inizializzazione
    _ "github.com/matteokutufa/zabbix-agent2-plugin-aws/aws"
    "github.com/matteokutufa/zabbix-agent2-plugin-aws/metrics/aws"
)

// Plugin è la struttura principale del plugin
type Plugin struct {
    plugin.Base
    options aws.Options
}

// impl è l'istanza del plugin
var impl Plugin

// Nome e versione del plugin
const (
	pluginName    = "AWS"
	pluginVersion = "1.0.0"
)

// Definizione delle opzioni di configurazione per il plugin
var (
	pluginComms       = comms.New("aws")
	configurationOpts = &conf.Opts{
		Listen: false,
	}

	// Parametri configurabili nel file di configurazione
	parameterOpts = []conf.ParameterOpts{
		{Key: "AccountFile", Type: conf.String, Default: "/etc/zabbix/aws_accounts.ini"},
		{Key: "MetricsFile", Type: conf.String, Default: "/etc/zabbix/metrics_config.yaml"},
		{Key: "Timeout", Type: conf.Int, Default: 30},
		{Key: "KeepAlive", Type: conf.Int, Default: 300},
		{Key: "Sessions", Type: conf.Int, Default: 15},
	}
)

// Configure configura il plugin con le opzioni specificate
func (p *Plugin) Configure(global *plugin.GlobalOptions, options interface{}) {
	if err := conf.Unmarshal(options, &p.options); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err)
	}

	aws.Configure(&p.options)
}

// Validate verifica le opzioni di configurazione
func (p *Plugin) Validate(options interface{}) error {
	var opts aws.Options

	return conf.Unmarshal(options, &opts)
}

// Export esporta le funzioni del plugin che possono essere chiamate da Zabbix
func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	switch key {
	// RDS
	case "aws.rds.discovery", "aws.rds.discovery.raw":
		return aws.RDSDiscovery(ctx, params, true)
	case "aws.rds.get", "aws.rds.get.raw":
		return aws.RDSGet(ctx, params, true)
	case "aws.rds.bulk", "aws.rds.bulk.raw":
		return aws.RDSBulkGet(ctx, params, true)

	// S3
	case "aws.s3.discovery", "aws.s3.discovery.raw":
		return aws.S3Discovery(ctx, params, true)
	case "aws.s3.get", "aws.s3.get.raw":
		return aws.S3Get(ctx, params, true)
	case "aws.s3.bulk", "aws.s3.bulk.raw":
		return aws.S3BulkGet(ctx, params, true)

	// ELB
	case "aws.elb.discovery", "aws.elb.discovery.raw":
		return aws.ELBDiscovery(ctx, params, true)
	case "aws.elb.get", "aws.elb.get.raw":
		return aws.ELBGet(ctx, params, true)
	case "aws.elb.bulk", "aws.elb.bulk.raw":
		return aws.ELBBulkGet(ctx, params, true)

	// ECS
	case "aws.ecs.cluster.discovery", "aws.ecs.cluster.discovery.raw":
		return aws.ECSClusterDiscovery(ctx, params, true)
	case "aws.ecs.service.discovery", "aws.ecs.service.discovery.raw":
		return aws.ECSServiceDiscovery(ctx, params, true)
	case "aws.ecs.get", "aws.ecs.get.raw":
		return aws.ECSGet(ctx, params, true)
	case "aws.ecs.bulk", "aws.ecs.bulk.raw":
		return aws.ECSBulkGet(ctx, params, true)

	// MSK
	case "aws.msk.discovery", "aws.msk.discovery.raw":
		return aws.MSKDiscovery(ctx, params, true)
	case "aws.msk.get", "aws.msk.get.raw":
		return aws.MSKGet(ctx, params, true)
	case "aws.msk.bulk", "aws.msk.bulk.raw":
		return aws.MSKBulkGet(ctx, params, true)

	// Ping - controllo di salute
	case "aws.ping":
		return aws.Ping(ctx, params, true)

	default:
		return nil, zbxerr.ErrorUnsupportedMetric
	}
}

func init() {
	plugin.RegisterMetrics(&impl, pluginName, pluginName+".ping",
		"aws.rds.discovery", "aws.rds.discovery.raw", "aws.rds.get", "aws.rds.get.raw", "aws.rds.bulk", "aws.rds.bulk.raw",
		"aws.s3.discovery", "aws.s3.discovery.raw", "aws.s3.get", "aws.s3.get.raw", "aws.s3.bulk", "aws.s3.bulk.raw",
		"aws.elb.discovery", "aws.elb.discovery.raw", "aws.elb.get", "aws.elb.get.raw", "aws.elb.bulk", "aws.elb.bulk.raw",
		"aws.ecs.cluster.discovery", "aws.ecs.cluster.discovery.raw",
		"aws.ecs.service.discovery", "aws.ecs.service.discovery.raw",
		"aws.ecs.get", "aws.ecs.get.raw", "aws.ecs.bulk", "aws.ecs.bulk.raw",
		"aws.msk.discovery", "aws.msk.discovery.raw", "aws.msk.get", "aws.msk.get.raw", "aws.msk.bulk", "aws.msk.bulk.raw")
}

func main() {
	// Inizializza il framework del plugin
	options := &plugin.Options{
		Name:                 pluginName,
		VersionCurrent:       pluginVersion,
		ConfigurationOptions: configurationOpts,
		ParameterOptions:     parameterOpts,
	}

	// Esegui il plugin
	plugin.Run(options, pluginComms)
}