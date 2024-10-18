package conf

// RabbitMQ contains all configuration related to rabbitmq.
type RabbitMQ struct {
	// Enable enables rabbitmq metrics and stats in the reports.
	Enable bool `koanf:"enable"`
	// Management contains configuration to access the rabbitmq management api.
	//
	// Required.
	Management RabbitMQManagement `kaonf:"management"`
	// HeadlessSvcAddr is the Kubernetes headless address pointing to
	// rabbitmq nodes. On a DNS lookup, this address must resolve to
	// rabbitmq node ips.
	// For example: rabbitmq-nodes.default.svc.cluster.local
	//
	// Required.
	HeadlessSvcAddr string `koanf:"headlessSvcAddr"`
	// Alerts contain a message template, a severity level, and a conditional
	// expression to trigger the respective alert.
	Alerts []Alert `koanf:"alerts"`
}

// RabbitMQManagement contains configuration to access the rabbitmq management
// api.
type RabbitMQManagement struct {
	// URL is the rabbitmq management url.
	// For example: http://rabbitmq.default.svc.cluster.local:15672
	//
	// Required.
	URL string `koanf:"url"`
	// Username is the basic auth username credential for the management api.
	//
	// Required.
	Username string `koanf:"username"`
	// Password is the basic auth password credential for the management api.
	//
	// Required.
	Password string `koanf:"password"`
}
