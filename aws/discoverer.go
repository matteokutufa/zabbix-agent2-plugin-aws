// File: aws/discoverer.go
package aws

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/rds"
    "github.com/matteokutufa/zabbix-agent2-plugin-aws/models"
)

// RDSDiscoverer è un'implementazione di RDSDiscovererInterface
type RDSDiscoverer struct {
    client models.AWSClientInterface
}

// NewRDSDiscoverer crea un nuovo discoverer RDS
func NewRDSDiscoverer(client models.AWSClientInterface) *RDSDiscoverer {
    return &RDSDiscoverer{
        client: client,
    }
}

// DiscoverInstances implementa il metodo dell'interfaccia
func (d *RDSDiscoverer) DiscoverInstances() (interface{}, error) {
    // Ottieni il client RDS
    rdsClient := d.client.RDSClient()

    // Crea l'input per la richiesta DescribeDBInstances
    input := &rds.DescribeDBInstancesInput{}

    // Esegui la chiamata API
    output, err := rdsClient.DescribeDBInstances(input)
    if err != nil {
        return nil, err
    }

    // Crea un slice per i risultati
    type dbInstance struct {
        DBInstance string `json:"{#DBINSTANCE}"`
        DBEngine   string `json:"{#DBENGINE}"`
        DBClass    string `json:"{#DBCLASS}"`
        DBStatus   string `json:"{#DBSTATUS}"`
    }

    type discoveryData struct {
        Data []dbInstance `json:"data"`
    }

    result := discoveryData{
        Data: make([]dbInstance, 0, len(output.DBInstances)),
    }

    // Itera su tutte le istanze RDS
    for _, instance := range output.DBInstances {
        result.Data = append(result.Data, dbInstance{
            DBInstance: aws.StringValue(instance.DBInstanceIdentifier),
            DBEngine:   aws.StringValue(instance.Engine),
            DBClass:    aws.StringValue(instance.DBInstanceClass),
            DBStatus:   aws.StringValue(instance.DBInstanceStatus),
        })
    }

    return result, nil
}