package configs

import "os"

type Config struct {
	UserName       string
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	Dev            bool
}

func LoadConfig(dev bool) *Config {
	return &Config{
		UserName:       "MyDailyBibleBot",
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		AccessToken:    os.Getenv("ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("ACCESS_SECRET"),
		Dev:            dev,
	}
}
