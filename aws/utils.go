// File: metrics/aws/utils.go (nuovo)
package aws

// MetricsFile restituisce il percorso del file delle metriche
func MetricsFile() string {
	return globalOptions.MetricsFile
}