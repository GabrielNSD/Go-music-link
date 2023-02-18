package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	DBDriver            string `mapstructure:"DB_DRIVER"`
	DBHOST              string `mapstructure:"DB_HOST"`
	JWTSecret           string `mapstructure:"DB_JWT_SECRET"`
	SpotifyClientID     string `mapstructure:"SPOTIFY_CLIENT_ID"`
	SpotifyClientSecret string `mapstructure:"SPOTIFY_CLIENT_SECRET"`
	TokenAuth           *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.TokenAuth = jwtauth.New("HS265", []byte(cfg.JWTSecret), nil)
	return cfg, err
}
