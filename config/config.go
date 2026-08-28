package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

func Load() error {
	v := viper.GetViper()
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	v.SetDefault("APP_PORT", "3000")
	v.SetDefault("AUTH_USERNAME", "admin")
	v.SetDefault("AUTH_PASSWORD", "admin")
	v.SetDefault("TOKEN_TTL_HOURS", 24)
	v.SetDefault("RATE_LIMIT_MAX", 5)
	v.SetDefault("RATE_LIMIT_WINDOW_MINUTES", 1)

	if err := v.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read config file: %w", err)
	}

	return nil
}

func Port() string {
	return viper.GetString("APP_PORT")
}

func Username() string {
	return viper.GetString("AUTH_USERNAME")
}

func Password() string {
	return viper.GetString("AUTH_PASSWORD")
}

func TokenTTL() time.Duration {
	return time.Duration(viper.GetInt("TOKEN_TTL_HOURS")) * time.Hour
}

func RateLimitMax() int {
	return viper.GetInt("RATE_LIMIT_MAX")
}

func RateLimitWindow() time.Duration {
	return time.Duration(viper.GetInt("RATE_LIMIT_WINDOW_MINUTES")) * time.Minute
}
