package config

import "github.com/spf13/viper"

type Config struct {
	Database          Database
	Scraper           Scraper
	Selectors         map[string]Selector
	CssContentExclude []string
}

type Database struct {
	Name       string
	Collection string
}

type Scraper struct {
	UserAgent string
	Timeout   int
}

type Selector struct {
	URL              string
	ArticleContainer string
	Title            string
	PublishedDate    DateSelector
	Content          string
	PageIndex        string
	Jsonld           string
}

type DateSelector struct {
	Css        string
	Attr       string
	TimeFormat string
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
