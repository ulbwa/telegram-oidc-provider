package config

type Config struct {
	HTTPServer HTTPServerConfig `yaml:"http_server" validate:"required"`
	Database   DatabaseConfig   `yaml:"database"    validate:"required"`
	Redis      RedisConfig      `yaml:"redis"       validate:"required"`
	Security   SecurityConfig   `yaml:"security"    validate:"required"`
	Hydra      HydraConfig      `yaml:"hydra"       validate:"required"`
	Logger     LoggerConfig     `yaml:"logger"      validate:"required"`
}
