package config

import "github.com/spf13/viper"

type Config struct {
	MongoDBName       string
	MongoCollection   string
	Sources           map[string]SourceConfig
	Selectors         map[string]SelectorConfig
	CssContentExclude []string
}

type SourceConfig struct {
	Domain   string
	StartURL string
}

type SelectorConfig struct {
	URL              string
	ArticleContainer string
	Title            string
	PublishedDate    struct {
		Css        string
		Attr       string
		TimeFormat string
	}
	Content   string
	PageIndex string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
