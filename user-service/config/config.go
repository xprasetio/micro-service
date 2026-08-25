package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `mapstructure:"APP_PORT" json:"app_port" yaml:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV" json:"app_env" yaml:"APP_ENV"`

	JwtSecretKey string `mapstructure:"JWT_SECRET_KEY" json:"jwt_secret_key" yaml:"JWT_SECRET_KEY"`
	JwtIssuer    string `mapstructure:"JWT_ISSUER" json:"jwt_issuer" yaml:"JWT_ISSUER"`

	UrlForgotPassword string `mapstructure:"URL_FORGOT_PASSWORD" json:"url_forgot_password" yaml:"URL_FORGOT_PASSWORD"`
}

type PsqlDB struct {
	Host      string `mapstructure:"PSQL_HOST" json:"psql_host" yaml:"PSQL_HOST"`
	Port      string `mapstructure:"PSQL_PORT" json:"psql_port" yaml:"PSQL_PORT"`
	User      string `mapstructure:"PSQL_USER" json:"psql_user" yaml:"PSQL_USER"`
	Password  string `mapstructure:"PSQL_PASSWORD" json:"psql_password" yaml:"PSQL_PASSWORD"`
	DBName    string `mapstructure:"PSQL_DB_NAME" json:"psql_db_name" yaml:"PSQL_DB_NAME"`
	SSLMode   string `mapstructure:"PSQL_SSL_MODE" json:"psql_ssl_mode" yaml:"PSQL_SSL_MODE"`
	DBMaxOpen int    `mapstructure:"PSQL_DB_MAX_OPEN" json:"psql_db_max_open" yaml:"PSQL_DB_MAX_OPEN"`
	DBMaxIdle int    `mapstructure:"PSQL_DB_MAX_IDLE" json:"psql_db_max_idle" yaml:"PSQL_DB_MAX_IDLE"`
}

type RabbitMQ struct {
	Host     string `mapstructure:"RABBITMQ_HOST" json:"rabbitmq_host" yaml:"RABBITMQ_HOST"`
	Port     string `mapstructure:"RABBITMQ_PORT" json:"rabbitmq_port" yaml:"RABBITMQ_PORT"`
	User     string `mapstructure:"RABBITMQ_USER" json:"rabbitmq_user" yaml:"RABBITMQ_USER"`
	Password string `mapstructure:"RABBITMQ_PASSWORD" json:"rabbitmq_password" yaml:"RABBITMQ_PASSWORD"`
}
type Config struct {
	App  App    `mapstructure:"app" json:"app" yaml:"app"`
	Psql PsqlDB `mapstructure:"psql" json:"psql" yaml:"psql"`

	RabbitMQ RabbitMQ `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort:      viper.GetString("APP_PORT"),
			AppEnv:       viper.GetString("APP_ENV"),
			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),

			UrlForgotPassword: viper.GetString("URL_FORGOT_PASSWORD"),
		},
		Psql: PsqlDB{
			Host:      viper.GetString("DATABASE_HOST"),
			Port:      viper.GetString("DATABASE_PORT"),
			User:      viper.GetString("DATABASE_USER"),
			Password:  viper.GetString("DATABASE_PASSWORD"),
			DBName:    viper.GetString("DATABASE_NAME"),
			SSLMode:   viper.GetString("DATABASE_SSL_MODE"),
			DBMaxOpen: viper.GetInt("DATABASE_MAX_OPEN_CONNECTIONS"),
			DBMaxIdle: viper.GetInt("DATABASE_MAX_IDLE_CONNECTIONS"),
		},

		RabbitMQ: RabbitMQ{
			Host:     viper.GetString("RABBITMQ_HOST"),
			Port:     viper.GetString("RABBITMQ_PORT"),
			User:     viper.GetString("RABBITMQ_USER"),
			Password: viper.GetString("RABBITMQ_PASSWORD"),
		},
	}
}
