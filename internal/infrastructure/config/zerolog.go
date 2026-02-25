package config

// LoggerConfig конфигурация логирования
type LoggerConfig struct {
	Level      string           `yaml:"level"`       // Глобальный уровень логирования (trace, debug, info, warn, error, fatal, panic)
	TimeFormat string           `yaml:"time_format"` // Формат времени (unix, unixms, unixmicro, rfc3339, rfc3339nano)
	Console    ConsoleLogConfig `yaml:"console"`     // Настройки вывода в консоль
	Files      []FileLogConfig  `yaml:"files"`       // Список файловых хендлеров
}

// ConsoleLogConfig конфигурация консольного вывода логов
type ConsoleLogConfig struct {
	Enabled  bool   `yaml:"enabled"`   // Включить/выключить вывод в консоль
	Level    string `yaml:"level"`     // Минимальный уровень логов (если пусто, используется глобальный)
	MaxLevel string `yaml:"max_level"` // Максимальный уровень логов (опционально, для фильтрации)
	Colored  bool   `yaml:"colored"`   // Использовать цветной вывод
	Pretty   bool   `yaml:"pretty"`    // Форматировать логи в читаемом виде (не JSON)
}

// FileLogConfig конфигурация файлового вывода логов
type FileLogConfig struct {
	Path     string       `yaml:"path"`      // Путь к файлу логов
	Level    string       `yaml:"level"`     // Минимальный уровень логов (если пусто, используется глобальный)
	MaxLevel string       `yaml:"max_level"` // Максимальный уровень логов (опционально, для фильтрации)
	Rotate   RotateConfig `yaml:"rotate"`    // Настройки ротации логов
}

// RotateConfig конфигурация ротации лог-файлов
type RotateConfig struct {
	Enabled    bool `yaml:"enabled"`     // Включить ротацию логов
	MaxSize    int  `yaml:"max_size"`    // Максимальный размер файла в мегабайтах
	MaxAge     int  `yaml:"max_age"`     // Максимальный возраст файлов в днях
	MaxBackups int  `yaml:"max_backups"` // Максимальное количество старых файлов
	Compress   bool `yaml:"compress"`    // Сжимать старые файлы
}
