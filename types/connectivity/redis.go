package connectivity

type Redis struct {
	Reachable bool `bson:"connected"`
}
