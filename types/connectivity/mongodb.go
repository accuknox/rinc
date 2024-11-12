package connectivity

type Mongodb struct {
	Reachable bool `bson:"connected"`
}
