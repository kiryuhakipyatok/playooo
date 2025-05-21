package config

import (
//	"fmt"
	"github.com/spf13/viper"
	"log"
	//"strings"
	//"github.com/caarlos0/env/v11"
)

type Config struct{
	App AppCfg
	Server ServerCfg
	Postgres PostgresCfg
	Redis RedisCfg
	Bot BotCfg
	Auth AuthCfg
}

type AppCfg struct{
	Name string `env:"APP_NAME,required"`
	Version string `env:"APP_VERSION,required"`
}

type ServerCfg struct{
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

type PostgresCfg struct{
	Host string `env:"PGHOST,required"`
	Port string `env:"PGPORT,required"`
	Database string `env:"POSTGRES_DB,required"`
	User string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
}

type RedisCfg struct{
	Host string `env:"REDISHOST,required"`
	Port string `env:"REDISPORT,required"`
	Password string `env:"REDIS_PASSWORD,required"`
}

type BotCfg struct{
	Token string `env:"TG_BOT_TOKEN,required"`
}

type AuthCfg struct{
	Secret string `env:"SECRET,required"`
}


// func LoadConfig() (*Config, error) {
// 	cfg := Config{}
// 	if err := env.Parse(&cfg); err != nil {
// 		return nil, err
// 	}

// 	return &cfg, nil
// }
func LoadConfig() (*Config,error){
	v:=viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("../../config")  
	v.AutomaticEnv()
	if err:=v.ReadInConfig();err!=nil{
		log.Printf("error reading config file: %v\n",err)
		return nil,err
	}
	cfg:=Config{}
	if err:=v.Unmarshal(&cfg);err!=nil{
		log.Printf("error unmarshaling config file: %v\n",err)
		return nil,err
	}

	return &cfg,nil
}