package config

import (
	"github.com/spf13/viper"
)

type App struct {
	Port        int    `mapstructure:"PORT"`
	Env         string `mapstructure:"ENV"`
	MaxFileSize int64  `mapstructure:"MAX_FILE_SIZE"`
}

type Service struct {
	Auth    string `mapstructure:"AUTH"`
	Backend string `mapstructure:"BACKEND"`
	File    string `mapstructure:"FILE"`
}

type Config struct {
	App     App
	Service Service
}

func LoadConfig() (*Config, error) {
	appCfgLdr := viper.New()
	appCfgLdr.SetEnvPrefix("APP")
	appCfgLdr.AutomaticEnv()
	appCfgLdr.AllowEmptyEnv(false)
	appConfig := App{}
	if err := appCfgLdr.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	serviceCfgLdr := viper.New()
	serviceCfgLdr.SetEnvPrefix("SERVICE")
	serviceCfgLdr.AutomaticEnv()
	serviceCfgLdr.AllowEmptyEnv(false)
	serviceConfig := Service{}
	if err := serviceCfgLdr.Unmarshal(&serviceConfig); err != nil {
		return nil, err
	}

	config := &Config{
		App:     appConfig,
		Service: serviceConfig,
	}

	return config, nil
}

func (ac *App) IsDevelopment() bool {
	return ac.Env == "development"
}
