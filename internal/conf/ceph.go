package conf

// Ceph contains all configuration related to ceph status reporter.
type Ceph struct {
	// Enable enables ceph status reporter.
	Enable bool `koanf:"enable"`
	// DashboardAPI contains configuration to access the ceph dashboard api.
	//
	// Required.
	DashboardAPI CephDashboardAPI `kaonf:"dashboardAPI"`
}

// CephDashboardAPI contains configuration to access the ceph dashboard API.
type CephDashboardAPI struct {
	// URL is the ceph dashboard API url.
	//
	// For example:
	// https://rook-ceph-mgr-dashboard.rook-ceph.svc.cluster.local:8443
	//
	// Required.
	URL string `koanf:"url"`
	// Username to authenticate with ceph dashboard API.
	//
	// Required.
	Username string `koanf:"username"`
	// Password to authenticate with ceph dashboard API.
	//
	// Required.
	Password string `koanf:"password"`
}
