// cloudwatch.go
package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// CloudWatchMetric rappresenta una metrica CloudWatch
type CloudWatchMetric struct {
	MetricName string
	Statistics string
	Dimensions map[string]string
	Period     int64
	StartTime  time.Time
	EndTime    time.Time
}

// GetRDSMetric recupera il valore di una metrica CloudWatch per un'istanza RDS specifica
func (c *Client) GetRDSMetric(instanceID, metricName, statistic string, period int64, startTime, endTime time.Time) (float64, error) {
	// Crea le dimensioni per la metrica RDS
	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("DBInstanceIdentifier"),
			Value: aws.String(instanceID),
		},
	}

	// Crea l'input per la richiesta GetMetricStatistics
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/RDS"),
		MetricName: aws.String(metricName),
		Dimensions: dimensions,
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Period:     aws.Int64(period),
		Statistics: []*string{
			aws.String(statistic),
		},
	}

	// Esegui la chiamata API
	output, err := c.cloudWatchClient.GetMetricStatistics(input)
	if err != nil {
		return 0, err
	}

	// Se non ci sono punti dati, restituisci 0
	if len(output.Datapoints) == 0 {
		return 0, nil
	}

	// Trova il datapoint più recente
	var latestDatapoint *cloudwatch.Datapoint

	for _, dp := range output.Datapoints {
		if latestDatapoint == nil || aws.TimeValue(dp.Timestamp).After(aws.TimeValue(latestDatapoint.Timestamp)) {
			latestDatapoint = dp
		}
	}

	// Estrai il valore in base alla statistica
	var value float64

	switch statistic {
	case "Average":
		value = aws.Float64Value(latestDatapoint.Average)
	case "Maximum":
		value = aws.Float64Value(latestDatapoint.Maximum)
	case "Minimum":
		value = aws.Float64Value(latestDatapoint.Minimum)
	case "Sum":
		value = aws.Float64Value(latestDatapoint.Sum)
	case "SampleCount":
		value = aws.Float64Value(latestDatapoint.SampleCount)
	default:
		// Se la statistica non è riconosciuta, utilizza Average
		value = aws.Float64Value(latestDatapoint.Average)
	}

	return value, nil
}

// ListAvailableRDSMetrics elenca tutte le metriche disponibili per un'istanza RDS
func (c *Client) ListAvailableRDSMetrics(instanceID string) ([]*cloudwatch.Metric, error) {
	// Crea l'input per la richiesta ListMetrics
	input := &cloudwatch.ListMetricsInput{
		Namespace: aws.String("AWS/RDS"),
		Dimensions: []*cloudwatch.DimensionFilter{
			{
				Name:  aws.String("DBInstanceIdentifier"),
				Value: aws.String(instanceID),
			},
		},
	}

	// Esegui la chiamata API
	output, err := c.cloudWatchClient.ListMetrics(input)
	if err != nil {
		return nil, err
	}

	return output.Metrics, nil
}
