package conf

import (
	"fmt"
	"os"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

var defaultConfig = "/etc/rinc/config.yaml"

// C contains all configuration data that can be passed to the reporter.
type C struct {
	RunAsGenerator bool
	RunAsWebServer bool
	// Log contains configuration for logs.
	Log Log `koanf:"log"`
	// TerminationGracePeriod is the period after which the web server
	// must be forcefully terminated. A value of 0 implies no forceful
	// termination.
	TerminationGracePeriod time.Duration `koanf:"terminationGracePeriod"`
	// Output is the path to the reports output directory.
	Output string `koanf:"output"`
	// KubernetesClient contains the configuration needed to communicate with
	// the Kubernetes API server.
	KubernetesClient KubernetesClient `koanf:"kubernetesClient"`
	// RabbitMQ contains the rabbitmq configuration.
	RabbitMQ RabbitMQ `koanf:"rabbitmq"`
	// LongJobs contains configuration related to the long-running job
	// reporter.
	LongJobs LongJobs `koanf:"longRunningJobs"`
}

// New creates a configuration using the provided arguments and config file.
func New(args ...string) (*C, error) {
	k := koanf.New(".")

	err := k.Load(confmap.Provider(map[string]any{
		"log.level":                        "info",
		"log.format":                       "text",
		"output":                           "reports",
		"terminationGracePeriod":           time.Second * 10,
		"kubernetesClient.inCluster":       false,
		"kubernetesClient.kubeconfig":      "",
		"rabbitmq.enable":                  false,
		"rabbitmq.headlessSvcAddr":         "",
		"rabbitmq.management.url":          "",
		"rabbitmq.management.username":     "",
		"rabbitmq.management.password":     "",
		"longRunningJobs.enable":           false,
		"longRunningJobs.namespace":        "ALL",
		"longRunningJobs.olderThan":        time.Hour * 12,
		"longRunningJobs.includeSuspended": false,
	}, "."), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load default configuration: %w", err)
	}

	f := parseFlags(args)
	confF, err := f.GetString("conf")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	asGenerator, err := f.GetBool("generate-reports")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	asWebServer, err := f.GetBool("serve")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	err = k.Load(file.Provider(confF), yaml.Parser())
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	conf := new(C)
	err = k.Unmarshal("", conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	conf.RunAsGenerator = asGenerator
	conf.RunAsWebServer = asWebServer

	return conf, nil
}

func parseFlags(args []string) *flag.FlagSet {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Print(f.FlagUsages())
		os.Exit(0)
	}
	f.String("conf", defaultConfig, "path to config file")
	f.Bool("generate-reports", false, "generate reports")
	f.Bool("serve", false, "serve static reports")
	f.Parse(args)
	return f
}
