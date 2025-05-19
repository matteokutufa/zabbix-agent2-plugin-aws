// client.go asdasd
package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/kafka"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/matteokutufa/zabbix-agent2-plugin-aws/aws/config"
)

// Client rappresenta un client AWS
type Client struct {
	session          *session.Session
	rdsClient        *rds.RDS
	cloudWatchClient *cloudwatch.CloudWatch
	s3Client         *s3.S3
	elbv2Client      *elbv2.ELBV2
	ecsClient        *ecs.ECS
	kafkaClient      *kafka.Kafka
	region           string
}

// NewClient crea un nuovo client AWS utilizzando le credenziali dell'account
func NewClient(account config.AWSAccount) (*Client, error) {
	// Crea una configurazione AWS
	awsConfig := &aws.Config{
		Region: aws.String(account.Region),
	}

	// Crea una sessione AWS
	var sess *session.Session
	var err error

	// Se è specificato un ruolo IAM, utilizziamo l'assunzione di ruolo
	if account.RoleARN != "" {
		// Crea una sessione base
		baseSess, err := session.NewSession(awsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create base session: %v", err)
		}

		// Crea le credenziali con l'assunzione di ruolo
		creds := stscreds.NewCredentials(baseSess, account.RoleARN)

		// Crea una nuova configurazione con le credenziali del ruolo
		awsConfig.Credentials = creds

		// Crea una sessione con le credenziali del ruolo
		sess, err = session.NewSession(awsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create session with role assumption: %v", err)
		}
	} else {
		// Altrimenti, utilizziamo le chiavi di accesso
		awsConfig.Credentials = credentials.NewStaticCredentials(
			account.AccessKey,
			account.SecretAccessKey,
			"", // Token di sessione, non necessario per le chiavi di accesso statiche
		)

		// Crea una sessione con le chiavi di accesso
		sess, err = session.NewSession(awsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create session with access keys: %v", err)
		}
	}

	// Crea i client per i servizi AWS
	rdsClient := rds.New(sess)
	cloudWatchClient := cloudwatch.New(sess)
	s3Client := s3.New(sess)
	elbv2Client := elbv2.New(sess)
	ecsClient := ecs.New(sess)
	kafkaClient := kafka.New(sess)

	return &Client{
		session:          sess,
		rdsClient:        rdsClient,
		cloudWatchClient: cloudWatchClient,
		s3Client:         s3Client,
		elbv2Client:      elbv2Client,
		ecsClient:        ecsClient,
		kafkaClient:      kafkaClient,
		region:           account.Region,
	}, nil
}

// RDSClient restituisce il client RDS
func (c *Client) RDSClient() *rds.RDS {
	return c.rdsClient
}

// CloudWatchClient restituisce il client CloudWatch
func (c *Client) CloudWatchClient() *cloudwatch.CloudWatch {
	return c.cloudWatchClient
}

// S3Client restituisce il client S3
func (c *Client) S3Client() *s3.S3 {
	return c.s3Client
}

// ELBv2Client restituisce il client ELBv2 (per ALB e NLB)
func (c *Client) ELBv2Client() *elbv2.ELBV2 {
	return c.elbv2Client
}

// ECSClient restituisce il client ECS
func (c *Client) ECSClient() *ecs.ECS {
	return c.ecsClient
}

// KafkaClient restituisce il client Kafka (MSK)
func (c *Client) KafkaClient() *kafka.Kafka {
	return c.kafkaClient
}

// Region restituisce la regione configurata
func (c *Client) Region() string {
	return c.region
}
