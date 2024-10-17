package ceph

type Status struct {
	Health     Health     `json:"health,omitempty" bson:"health,omitempty"`
	OSDMap     OSDMap     `json:"osd_map,omitempty" bson:"osd_map,omitempty"`
	Pools      []Pool     `json:"pools,omitempty" bson:"pools,omitempty"`
	MGRMap     MGRMap     `json:"mgr_map,omitempty" bson:"mgr_map,omitempty"`
	PGInfo     PGInfo     `json:"pg_info,omitempty" bson:"pg_info,omitempty"`
	MonStatus  MonStatus  `json:"mon_status,omitempty" bson:"mon_status,omitempty"`
	DF         DF         `json:"df,omitempty" bson:"df,omitempty"`
	ClientPerf ClientPerf `json:"client_perf,omitempty" bson:"client_perf,omitempty"`
	Hosts      uint       `json:"hosts,omitempty" bson:"hosts,omitempty"`
}

type Health struct {
	Status string  `json:"status,omitempty" bson:"status,omitempty"`
	Checks []Check `json:"checks,omitempty" bson:"checks,omitempty"`
}

type Check struct {
	Severity string   `json:"severity,omitempty" bson:"severity,omitempty"`
	Detail   []Detail `json:"detail,omitempty" bson:"detail,omitempty"`
	Muted    bool     `json:"muted,omitempty" bson:"muted,omitempty"`
	Type     string   `json:"type,omitempty" bson:"type,omitempty"`
}

type Detail struct {
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type OSDMap struct {
	OSDs []OSD `json:"osds,omitempty" bson:"osds,omitempty"`
}

type OSD struct {
	Up    uint     `json:"up,omitempty" bson:"up,omitempty"`
	In    uint     `json:"in,omitempty" bson:"in,omitempty"`
	State []string `json:"state,omitempty" bson:"state,omitempty"`
}

type Pool struct {
	PGNum uint `json:"pg_num,omitempty" bson:"pg_num,omitempty"`
}

type MGRMap struct {
	ActiveName string   `json:"active_name,omitempty" bson:"active_name,omitempty"`
	StandBys   []string `json:"standbys,omitempty" bson:"standbys,omitempty"`
}

type PGInfo struct {
	Statuses  map[string]uint `json:"statuses,omitempty" bson:"statuses,omitempty"`
	PGsPerOSD float64         `json:"pgs_per_osd,omitempty" bson:"pgs_per_osd,omitempty"`
}

type MonStatus struct {
	MonMap MonMap `json:"monmap,omitempty" bson:"monmap,omitempty"`
}

type MonMap struct {
	Mon []Mon `json:"mons,omitempty" bson:"mons,omitempty"`
}

type Mon struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

type DF struct {
	Stats DFStats `json:"stats,omitempty" bson:"stats,omitempty"`
}

type DFStats struct {
	TotalBytes      uint64 `json:"total_bytes,omitempty" bson:"total_bytes,omitempty"`
	TotalAvailBytes uint64 `json:"total_avail_bytes,omitempty" bson:"total_avail_bytes,omitempty"`
	TotalUsedBytes  uint64 `json:"total_used_raw_bytes,omitempty" bson:"total_used_raw_bytes,omitempty"`
}

type ClientPerf struct {
	ReadBytesPerSec       uint64 `json:"read_bytes_sec,omitempty" bson:"read_bytes_sec,omitempty"`
	ReadOpPerSec          uint64 `json:"read_op_per_sec,omitempty" bson:"read_op_per_sec,omitempty"`
	WriteBytesPerSec      uint64 `json:"write_bytes_sec,omitempty" bson:"write_bytes_sec,omitempty"`
	WriteOpPerSec         uint64 `json:"write_op_per_sec,omitempty" bson:"write_op_per_sec,omitempty"`
	RecoveringBytesPerSec uint64 `json:"recovering_bytes_per_sec,omitempty" bson:"recovering_bytes_per_sec,omitempty"`
}
