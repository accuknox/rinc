package conf

type Mongodb struct {
	URI      string `koanf:"uri"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}
