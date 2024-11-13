package kube

import (
	"fmt"

	"github.com/accuknox/rinc/internal/conf"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// NewClient creates a new Kubernetes API server client using the provided
// configuration.
func NewClient(c conf.KubernetesClient) (*kubernetes.Clientset, error) {
	conf, err := config(c)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("creating new kube client: %w", err)
	}
	return client, nil
}

// NewMetricsClient creates a new Kubernetes Metrics API client using the
// provided configuration.
func NewMetricsClient(c conf.KubernetesClient) (*metrics.Clientset, error) {
	conf, err := config(c)
	if err != nil {
		return nil, err
	}
	client, err := metrics.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("creating new metrics client: %w", err)
	}
	return client, err
}

func config(c conf.KubernetesClient) (*rest.Config, error) {
	if c.InCluster {
		conf, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("in cluster config: %w", err)
		}
		return conf, nil
	}
	if c.Kubeconfig == "" {
		return nil, fmt.Errorf("either `InCluster` or `Kubeconfig` must be set")
	}
	conf, err := clientcmd.BuildConfigFromFlags("", c.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("config from kubeconfig: %w", err)
	}
	return conf, nil
}
