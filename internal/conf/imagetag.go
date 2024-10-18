package conf

// ImageTag contains configuration related to the image tag reporter.
type ImageTag struct {
	// Enable specifies whether the image tag reporter is enabled.
	Enable bool `koanf:"enable"`
	// Namespace is the Kubernetes namespace that the image tag reporter
	// will be limited to.
	Namespace string `koanf:"namespace"`
	// Alerts contain a message template, a severity level, and a conditional
	// expression to trigger the respective alert.
	Alerts []Alert `koanf:"alerts"`
}
