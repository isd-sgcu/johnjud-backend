package config

import (
	"github.com/spf13/viper"
)

type App struct {
	Port int    `mapstructure:"PORT"`
	Env  string `mapstructure:"ENV"`
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
	appConfig := App{}
	LoadEnvGroup(&appConfig, "APP")

	serviceConfig := Service{}
	LoadEnvGroup(&serviceConfig, "SERVICE")

	config := &Config{
		App:     appConfig,
		Service: serviceConfig,
	}

	return config, nil
}

func (ac *App) IsDevelopment() bool {
	return ac.Env == "development"
}

func LoadEnvGroup(config interface{}, prefix string) (err error) {
	cfgLdr := viper.New()
	cfgLdr.SetEnvPrefix(prefix)
	cfgLdr.AutomaticEnv()
	cfgLdr.AllowEmptyEnv(false)
	if err := cfgLdr.Unmarshal(&config); err != nil {
		return err
	}
	return nil
}
