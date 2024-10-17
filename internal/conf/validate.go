package conf

import (
	"fmt"
	"net"
)

// Validate validates the provided configuration.
func (c C) Validate() error {
	if err := validateLogLevel(c.Log.Level); err != nil {
		return fmt.Errorf("`log.level`: %w", err)
	}
	if err := validateLogFormat(c.Log.Format); err != nil {
		return fmt.Errorf("`log.format`: %w", err)
	}
	if err := validateKubernetesClient(c.KubernetesClient); err != nil {
		return fmt.Errorf("`kubernetesClient`: %w", err)
	}
	if err := validateMongodb(c.Mongodb); err != nil {
		return fmt.Errorf("`mongodb`: %w", err)
	}
	if err := validateRabbitMQ(c.RabbitMQ); err != nil {
		return fmt.Errorf("rabbitmq: %w", err)
	}
	if err := validateCeph(c.Ceph); err != nil {
		return fmt.Errorf("ceph: %w", err)
	}
	return nil
}

func validateLogLevel(level string) error {
	switch level {
	case "debug":
	case "info":
	case "warn":
	case "error":
	default:
		return fmt.Errorf("invalid value for `log.level`: %q", level)
	}
	return nil
}

func validateLogFormat(format string) error {
	switch format {
	case "text":
	case "json":
	default:
		return fmt.Errorf("invalid value for `log.format`: %q", format)
	}
	return nil
}

func validateKubernetesClient(c KubernetesClient) error {
	if c.InCluster || c.Kubeconfig != "" {
		return nil
	}
	return fmt.Errorf("either `inCluster` or `kubeconfig` must be set")
}

func validateMongodb(c Mongodb) error {
	if c.URI == "" {
		return fmt.Errorf("missing `mongodb.uri`")
	}
	if c.Username == "" {
		return fmt.Errorf("missing `mongodb.username`")
	}
	if c.Password == "" {
		return fmt.Errorf("missing `mongodb.password`")
	}
	return nil
}

func validateRabbitMQ(rmq RabbitMQ) error {
	if !rmq.Enable {
		return nil
	}
	if rmq.Management.URL == "" {
		return fmt.Errorf("missing `rabbitmq.management.url`")
	}
	if rmq.Management.Username == "" {
		return fmt.Errorf("missing `rabbitmq.management.username`")
	}
	if rmq.Management.Password == "" {
		return fmt.Errorf("missing `rabbitmq.management.password`")
	}
	if rmq.HeadlessSvcAddr == "" {
		return fmt.Errorf("missing `rabbitmq.headlessSvcAddr`")
	}
	_, err := net.LookupIP(rmq.HeadlessSvcAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve %q: %w", rmq.HeadlessSvcAddr, err)
	}
	return nil
}

func validateCeph(c Ceph) error {
	if !c.Enable {
		return nil
	}
	if c.DashboardAPI.URL == "" {
		return fmt.Errorf("missing `ceph.dashboardAPI.url`")
	}
	if c.DashboardAPI.Username == "" {
		return fmt.Errorf("missing `ceph.dashboardAPI.username`")
	}
	if c.DashboardAPI.Password == "" {
		return fmt.Errorf("missing `ceph.dashboardAPI.password`")
	}
	return nil
}
