package serverconfig

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	IsProduction bool `mapstructure:"release"`
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("release", false)
}

func New() (*Config, error) {
	v := viper.New()

	setDefaults(v)
	setEnv(v)
	errConfigFile := setConfigFile(v)

	var result Config
	err := v.Unmarshal(&result)
	if err != nil {
		return nil, fmt.Errorf("decoding into struct: %w", err)
	}

	return &result, errConfigFile
}

func setEnv(v *viper.Viper) {
	v.SetEnvPrefix("mess")
	v.AutomaticEnv()
}

func setConfigFile(v *viper.Viper) error {
	v.SetConfigName("mess-server")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return errors.Join(err, ErrConfigFileNotFound)
	}
	return nil
}

var ErrConfigFileNotFound = fmt.Errorf("")
