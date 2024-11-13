package conf

// Connectivity contains all configuration related to connectivity
// status reporter.
type Connectivity struct {
	// Vault contains all configuration related to vault connectivity check.
	Vault VaultCheck `koanf:"vault"`
	// Mongodb contains all configuration related to mongodb connectivity
	// check.
	Mongodb MongodbCheck `koanf:"mongodb"`
	// Neo4j contains all configuration related to neo4j connectivity check.
	Neo4j Neo4jCheck `koanf:"neo4j"`
	// Postgres contains all configuration related to postgres connectivity check.
	Postgres PostgresCheck `koanf:"postgres"`
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
	// URI is the mongodb connection uri.
	//
	// E.g., mongodb://accuknox-mongodb-rs0.accuknox-mongodb.svc.cluster.local:27017
	URI string `koanf:"uri"`
}

// Neo4j contains all configuration related to neo4j connectivity check.
type Neo4jCheck struct {
	// Enable enables neo4j connectivity check.
	Enable bool `koanf:"enable"`
	// URI is the neo4j connection uri.
	//
	// E.g., neo4j://neo4j.accuknox-neo4j.svc.cluster.local:7687
	URI string `koanf:"uri"`
	// Username is the neo4j basic auth username.
	Username string `koanf:"username"`
	// Password is the neo4j basic auth password.
	Password string `koanf:"password"`
}

// Postgres contains all configuration related to postgres connectivity check.
type PostgresCheck struct {
	// Enable enables postgres connectivity check.
	Enable bool `koanf:"enable"`
	// Host is the postgresql server host (without the port).
	//
	// E.g., postgres-replicas.accuknox-postgresql.svc.cluster.local
	Host string `koanf:"host"`
	// Port is the postgresql server port.
	//
	// Default: 5432
	Port uint16 `koanf:"port"`
	// Username is the postgres auth username.
	Username string `koanf:"username"`
	// Password is the postgres auth password.
	Password string `koanf:"password"`
}
