package ceph

type Status struct {
	Health     Health     `json:"health"`
	OSDMap     OSDMap     `json:"osd_map"`
	Pools      []Pool     `json:"pools"`
	MGRMap     MGRMap     `json:"mgr_map"`
	PGInfo     PGInfo     `json:"pg_info"`
	MonStatus  MonStatus  `json:"mon_status"`
	DF         DF         `json:"df"`
	ClientPerf ClientPerf `json:"client_perf"`
	Hosts      uint       `json:"hosts"`
}

type Health struct {
	Status string  `json:"status"`
	Checks []Check `json:"checks"`
}

type Check struct {
	Severity string   `json:"severity"`
	Detail   []Detail `json:"detail"`
	Muted    bool     `json:"muted"`
	Type     string   `json:"type"`
}

type Detail struct {
	Message string `json:"message"`
}

type OSDMap struct {
	OSDs []OSD `json:"osds"`
}

type OSD struct {
	Up    uint     `json:"up"`
	In    uint     `json:"in"`
	State []string `json:"state"`
}

type Pool struct {
	PGNum uint `json:"pg_num"`
}

type MGRMap struct {
	ActiveName string   `json:"active_name"`
	StandBys   []string `json:"standbys"`
}

type PGInfo struct {
	Statuses  map[string]uint `json:"statuses"`
	PGsPerOSD float64         `json:"pgs_per_osd"`
}

type MonStatus struct {
	MonMap MonMap `json:"monmap"`
}

type MonMap struct {
	Mon []Mon `json:"mons"`
}

type Mon struct {
	Name string `json:"name"`
}

type DF struct {
	Stats DFStats `json:"stats"`
}

type DFStats struct {
	TotalBytes      uint64 `json:"total_bytes"`
	TotalAvailBytes uint64 `json:"total_avail_bytes"`
	TotalUsedBytes  uint64 `json:"total_used_raw_bytes"`
}

type ClientPerf struct {
	ReadBytesPerSec       uint64 `json:"read_bytes_sec"`
	ReadOpPerSec          uint64 `json:"read_op_per_sec"`
	WriteBytesPerSec      uint64 `json:"write_bytes_sec"`
	WriteOpPerSec         uint64 `json:"write_op_per_sec"`
	RecoveringBytesPerSec uint64 `json:"recovering_bytes_per_sec"`
}
