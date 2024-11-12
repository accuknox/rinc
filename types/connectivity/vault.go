package connectivity

type Vault struct {
	Reachable   bool   `bson:"connected"`
	Initialized bool   `bson:"initialized"`
	Sealed      bool   `bson:"sealed"`
	Version     string `bson:"version"`
	ClusterName string `bson:"clusterName"`
}
