package conf

// Connectivity contains all configuration related to connectivity
// status reporter.
type Connectivity struct {
	// Vault contains all configuration related to vault connectivity check.
	Vault Vault `koanf:"vault"`
	// Alerts contain a message template, a severity level, and a
	// conditional expression to trigger the respective alert.
	Alerts []Alert `koanf:"alerts"`
}

// Vault contains all configuration related to vault connectivity check.
type Vault struct {
	// Enable enables vault connectivity check.
	Enable bool `koanf:"enable"`
	// Addr is the vault address.
	//
	// E.g., http://accuknox-vault.accuknox-vault.svc.cluster.local:8200
	Addr string `koanf:"addr"`
}
