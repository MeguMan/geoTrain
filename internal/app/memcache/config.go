package memcache

type Config struct {
	Password   string `json:"password"`
	Capacity   int    `json:"capacity"`
	SessionKey string `json:"session_key"`
}

func NewConfig() *Config {
	return &Config{}
}

