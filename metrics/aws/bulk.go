// File: metrics/aws/bulk.go (senza RDSBulkGet)
package aws

// MetricResult rappresenta il risultato di una metrica
type MetricResult struct {
	MetricName string  `json:"metric"`
	Value      float64 `json:"value"`
	Timestamp  int64   `json:"timestamp"`
	Statistic  string  `json:"statistic"`
}

// BulkResult rappresenta il risultato di una richiesta bulk
type BulkResult struct {
	ResourceID string         `json:"resource_id"`
	Metrics    []MetricResult `json:"metrics"`
}

// S3BulkGet ottiene tutte le metriche S3 in un'unica chiamata
//func S3BulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
//	// Implementazione simile per S3
//	// ...
//
//
// ELBBulkGet ottiene tutte le metriche ELB in un'unica chiamata
//func ELBBulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
//	// Implementazione simile per ELB
//	// ...
//
//
/// ECSBulkGet ottiene tutte le metriche ECS in un'unica chiamata
//func ECSBulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
//	// Implementazione simile per ECS
//	// ...
//
//
/// MSKBulkGet ottiene tutte le metriche MSK in un'unica chiamata
//func MSKBulkGet(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
//	// Implementazione simile per MSK
//	// ...
//
//