// CloudWatchMetric è stato spostato da aws/cloudwatch.go a models/types.go
package models

import (
	"time"
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