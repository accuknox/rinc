package conf

// ResourceUtilization contains configuration related to the resource
// utilization reporter.
type ResourceUtilization struct {
	// Enable specifies whether the resource utilization reporter is enabled.
	Enable bool `koanf:"enable"`
	// Namespace is the Kubernetes namespace that the resource utilization
	// reporter will be limited to.
	Namespace string `koanf:"namespace"`
	// Alerts contain a message template, a severity level, and a conditional
	// expression to trigger the respective alert.
	Alerts []Alert `koanf:"alerts"`
}
