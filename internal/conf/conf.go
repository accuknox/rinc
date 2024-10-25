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
	RunAsScraper   bool
	RunAsWebServer bool
	GenerateSchema string
	// Log contains configuration for logs.
	Log Log `koanf:"log"`
	// TerminationGracePeriod is the period after which the web server
	// must be forcefully terminated. A value of 0 implies no forceful
	// termination.
	TerminationGracePeriod time.Duration `koanf:"terminationGracePeriod"`
	// KubernetesClient contains the configuration needed to communicate with
	// the Kubernetes API server.
	KubernetesClient KubernetesClient `koanf:"kubernetesClient"`
	Mongodb          Mongodb          `koanf:"mongodb"`
	// RabbitMQ contains the rabbitmq configuration.
	RabbitMQ RabbitMQ `koanf:"rabbitmq"`
	// LongJobs contains configuration related to the long-running job
	// reporter.
	LongJobs LongJobs `koanf:"longRunningJobs"`
	// ImageTag contains configuration related to the image tag
	// reporter.
	ImageTag ImageTag `koanf:"imageTag"`
	// DaSS contains configuration related to the deployment and statefulset
	// status reporter.
	DaSS DaSS `koanf:"deploymentAndStatefulsetStatus"`
	// Ceph contains configuration related to the ceph status reporter.
	Ceph Ceph `koanf:"ceph"`
}

// New creates a configuration using the provided arguments and config file.
func New(args ...string) (*C, error) {
	k := koanf.New(".")

	err := k.Load(confmap.Provider(map[string]any{
		"log.level":                 "info",
		"log.format":                "text",
		"terminationGracePeriod":    time.Second * 10,
		"longRunningJobs.olderThan": time.Hour * 12,
	}, "."), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load default configuration: %w", err)
	}

	f := parseFlags(args)
	confF, err := f.GetStringSlice("conf")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	asScraper, err := f.GetBool("scrape")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	asWebServer, err := f.GetBool("serve")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	generateSchema, err := f.GetString("generate-schema")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	for _, c := range confF {
		err := k.Load(file.Provider(c), yaml.Parser())
		if err != nil {
			return nil, fmt.Errorf("failed to load config %q: %w", c, err)
		}
	}

	conf := new(C)
	err = k.Unmarshal("", conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	conf.RunAsScraper = asScraper
	conf.RunAsWebServer = asWebServer
	conf.GenerateSchema = generateSchema

	return conf, nil
}

func parseFlags(args []string) *flag.FlagSet {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Print(f.FlagUsages())
		os.Exit(0)
	}
	f.StringSlice("conf", []string{defaultConfig}, "comma-seperated list of config files")
	f.String("generate-schema", "", "generate json schema")
	f.Bool("scrape", false, "scrape & store metrics")
	f.Bool("serve", false, "serve static reports")
	f.Parse(args)
	return f
}
