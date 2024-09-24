package ceph

const (
	PermissionActionRead   = "read"
	PermissionActionWrite  = "write"
	PermissionActionCreate = "create"
	PermissionActionDelete = "delete"
)

const PermissionHosts = "hosts"

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	T           string              `json:"token"`
	Permissions map[string][]string `json:"permissions"`
}
