package config

import (
	"time"

	viper "github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf  `mapstructure:"logger"`
	Storage StorageConf `mapstructure:"storage"`
	SQL     SQLConfig   `mapstructure:"sql"`
	HTTP    HTTPConfig  `mapstructure:"http"`
}

type LoggerConf struct {
	Level     string `mapstructure:"level"`
	Format    string `mapstructure:"format"`
	AddSource bool   `mapstructure:"addsource"`
}
type StorageConf struct {
	Type string `mapstructure:"type"`
}
type SQLConfig struct {
	Dsn string `mapstructure:"dsn"`
}
type HTTPConfig struct {
	Address        string        `mapstructure:"addr"`
	RequestTimeout time.Duration `mapstructure:"requesttimeout"`
}

func NewConfig() Config {
	return Config{}
}

func (config *Config) ReadConfig(configFile string) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	return nil
}
