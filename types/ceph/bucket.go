package ceph

type Bucket struct {
	Name       string      `json:"bucket,omitempty" bson:"bucket,omitempty"`
	NumShards  uint        `json:"num_shards,omitempty" bson:"num_shards,omitempty"`
	Tenant     string      `json:"tenant,omitempty" bson:"tenant,omitempty"`
	ZoneGroup  string      `json:"zone_group,omitempty" bson:"zone_group,omitempty"`
	MFAEnabled bool        `json:"mfa_enabled,omitempty" bson:"mfa_enabled,omitempty"`
	Owner      string      `json:"owner,omitempty" bson:"owner,omitempty"`
	Quota      BucketQuota `json:"bucket_quota,omitempty" bson:"bucket_quota,omitempty"`
	Usage      BucketUsage `json:"usage,omitempty" bson:"usage,omitempty"`
}

type BucketQuota struct {
	Enabled    bool  `json:"enabled,omitempty" bson:"enabled,omitempty"`
	MaxSize    int64 `json:"max_size,omitempty" bson:"max_size,omitempty"`
	MaxObjects int64 `json:"max_objects,omitempty" bson:"max_objects,omitempty"`
}

type BucketUsage struct {
	Main      BucketUsageStats `json:"rgw.main,omitempty" bson:"rgw.main,omitempty"`
	Multimeta BucketUsageStats `json:"rgw.multimeta,omitempty" bson:"rgw.multimeta,omitempty"`
}

type BucketUsageStats struct {
	Size         uint64 `json:"size_actual,omitempty" bson:"size_actual,omitempty"`
	SizeUtilized uint64 `json:"size_utilized,omitempty" bson:"size_utilized,omitempty"`
	NumObjects   uint   `json:"num_objects,omitempty" bson:"num_objects,omitempty"`
}
