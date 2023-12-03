package serverconfig

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const defaultSessionSecretLen = 64

type Config struct {
	IsProduction       bool          `mapstructure:"release"`
	SessionSecret      string        `mapstructure:"session_secret"`
	Port               int           `mapstructure:"port"`
	IncomingOrigin     string        `mapstructure:"incoming_origin"`
	AssetsCacheMaxAge  int           `mapstructure:"assets_cache_max_age"`
	HeartbeatPeriod    time.Duration `mapstructure:"heartbeat_period"`
	MaxWebsocketErrors int           `mapstructure:"max_websocket_errors"`
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("release", false)
	v.SetDefault("session_secret", generateSessionSecret())
	v.SetDefault("port", 4000)
	v.SetDefault("incoming_origin", "http://localhost:3000")
	v.SetDefault("assets_cache_max_age", 600)
	v.SetDefault("heartbeat_period", time.Second*5)
	v.SetDefault("max_websocket_errors", 5)
}

func generateSessionSecret() string {
	result := make([]byte, defaultSessionSecretLen)
	n, err := rand.Read(result)
	if err != nil {
		panic(err)
	} else if n != defaultSessionSecretLen {
		err := fmt.Errorf(
			"can't generate session secret: read %v of %v bytes",
			n, defaultSessionSecretLen)
		panic(err)
	}
	return string(result)
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
