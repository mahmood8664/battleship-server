package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var (
	envPrefix = "BATTLESHIP"
	C         Config
)

type Config struct {
	Port    string  `yaml:"port"`
	Logging Logging `yaml:"logging"`
}

type Logging struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	FileName string `yaml:"file_name"`
}

func Init(filename string) *Config {
	c := Config{}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if filename != "" {
		viper.SetConfigFile(filename)
		if err := viper.MergeInConfig(); err != nil {
			logrus.Fatalf("loading configs file [%s] failed: %s", filename, err)
		} else {
			logrus.Infof("configs file [%s] loaded successfully", filename)
		}
	} else {
		logrus.Fatal("Config fil is not determined")
	}

	err := viper.Unmarshal(&c, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
	})
	if err != nil {
		logrus.Fatal("failed on configs unmarshal: ", err)
	}

	C = c
	logrus.Debugf("Following configuration is loaded:\n%+v\n", c)
	return &c
}
