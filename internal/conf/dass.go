package conf

// DaSS contains configuration related to the deployment and statefulset status
// reporter.
type DaSS struct {
	// Enable specifies whether the deployment and statefulset status (DaSS)
	// reporter is enabled.
	Enable bool `koanf:"enable"`
	// Namespace is the Kubernetes namespace that the DaSS reporter will be
	// limited to. Leave blank for all namespaces.
	Namespace string `koanf:"namespace"`
	// Alerts contain a message template, a severity level, and a conditional
	// expression to trigger the respective alert.
	Alerts []Alert `koanf:"alerts"`
}
