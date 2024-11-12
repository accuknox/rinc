package connectivity

import "time"

type Metrics struct {
	Timestamp time.Time `bson:"timestamp"`
	Vault     Vault     `bson:"vault"`
}
