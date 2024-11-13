package connectivity

import "time"

type Metrics struct {
	Timestamp time.Time `bson:"timestamp"`
	Vault     Vault     `bson:"vault"`
	Mongodb   Mongodb   `bson:"mongodb"`
	Neo4j     Neo4j     `bson:"neo4j"`
	Postgres  Postgres  `bson:"postgres"`
}
