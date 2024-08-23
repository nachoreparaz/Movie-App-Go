package config

import "github.com/spf13/viper"

type Config struct {
	ServerAddress string
	DatabaseURL   string
	MoviesbaseURL string
	API_KEY       string
	SecretJWT     string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerAddress: viper.GetString("server_address"),
		DatabaseURL:   viper.GetString("database_url"),
		MoviesbaseURL: viper.GetString("movies_base_url"),
		API_KEY:       viper.GetString("api_key"),
		SecretJWT:     viper.GetString("secret"),
	}

	return cfg, nil

}
