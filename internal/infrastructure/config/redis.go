package config

// RedisConfig represents Redis connection configuration.
type RedisConfig struct {
	Addr     string `yaml:"addr" validate:"required"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
