package conf

// Connectivity contains all configuration related to connectivity
// status reporter.
type Connectivity struct {
	// Vault contains all configuration related to vault connectivity check.
	Vault VaultCheck `koanf:"vault"`
	// Mongodb contains all configuration related to mongodb connectivity
	// check.
	Mongodb MongodbCheck `koanf:"mongodb"`
	// Alerts contain a message template, a severity level, and a
	// conditional expression to trigger the respective alert.
	Alerts []Alert `koanf:"alerts"`
}

// VaultCheck contains all configuration related to vault connectivity check.
type VaultCheck struct {
	// Enable enables vault connectivity check.
	Enable bool `koanf:"enable"`
	// Addr is the vault address.
	//
	// E.g., http://accuknox-vault.accuknox-vault.svc.cluster.local:8200
	Addr string `koanf:"addr"`
}

// Mongodb contains all configuration related to mongodb connectivity check.
type MongodbCheck struct {
	// Enable enables mongodb connectivity check.
	Enable bool `koanf:"enable"`
	// URI is the mongodb uri.
	//
	// E.g., mongodb://accuknox-mongodb-rs0.accuknox-mongodb.svc.cluster.local:27017
	URI string `koanf:"uri"`
}
