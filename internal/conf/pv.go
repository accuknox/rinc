package conf

// PVUtilization contains configuration related to the PV utilization reporter.
type PVUtilization struct {
	// Enable specifies whether the PV utilization reporter is enabled.
	Enable bool `koanf:"enable"`
	// PrometheusURL is the prometheus service url. PV utilization reporter
	// depend on Prometheus to fetch the utilization.
	//
	// E.g., http://prometheus.monitoring.svc.cluster.local:9090
	PrometheusURL string `koanf:"prometheusUrl"`
	// Alerts contain a message template, a severity level, and a conditional
	// expression to trigger the respective alert.
	Alerts []Alert `koanf:"alerts"`
}
