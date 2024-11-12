package conf

// PodStatus contains configuration related to the pod status reporter.
type PodStatus struct {
	// Enable specifies whether the pod status reporter is enabled.
	Enable bool `koanf:"enable"`
	// Namespace is the Kubernetes namespace that the pod status reporter will
	// be limited to. Leave blank for all namespaces.
	Namespace string `koanf:"namespace"`
	// Alerts contain a message template, a severity level, and a conditional
	// expression to trigger the respective alert.
	Alerts []Alert `koanf:"alerts"`
}
