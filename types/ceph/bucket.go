package ceph

type Bucket struct {
	Name       string      `json:"bucket"`
	NumShards  uint        `json:"num_shards"`
	Tenant     string      `json:"tenant"`
	ZoneGroup  string      `json:"zone_group"`
	MFAEnabled bool        `json:"mfa_enabled"`
	Owner      string      `json:"owner"`
	Quota      BucketQuota `json:"bucket_quota"`
	Usage      BucketUsage `json:"usage"`
}

type BucketQuota struct {
	Enabled    bool  `json:"enabled"`
	MaxSize    int64 `json:"max_size"`
	MaxObjects int64 `json:"max_objects"`
}

type BucketUsage struct {
	Main      BucketUsageStats `json:"rgw.main"`
	Multimeta BucketUsageStats `json:"rgw.multimeta"`
}

type BucketUsageStats struct {
	Size         uint64 `json:"size_actual"`
	SizeUtilized uint64 `json:"size_utilized"`
	NumObjects   uint   `json:"num_objects"`
}
