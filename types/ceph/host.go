package ceph

import "time"

type Host struct {
	Hostname string   `json:"hostname"`
	Status   string   `json:"status"`
	Addr     string   `json:"addr"`
	Labels   []string `json:"labels"`
}

type Device struct {
	ID       string           `json:"devid"`
	Location []DeviceLocation `json:"location"`
}

type Inventory struct {
	Hostname      string         `json:"name"`
	PhysicalDisks []PhysicalDisk `json:"devices"`
}

type DeviceLocation struct {
	Host string `json:"host"`
	Dev  string `json:"dev"`
	Path string `json:"path"`
}

type PhysicalDisk struct {
	RejectedReasons []string          `json:"rejected_reasons"`
	Available       bool              `json:"available"`
	Path            string            `json:"path"`
	Stats           PhysicalDiskStats `json:"sys_api"`
	Created         time.Time         `json:"created"`
	Type            string            `json:"human_readable_type"`
}

type PhysicalDiskStats struct {
	Size float64 `json:"size"`
}
