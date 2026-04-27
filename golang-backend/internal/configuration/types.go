package configuration

type DatabaseConfiguration struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

type ServerConfiguration struct {
	Port int `mapstructure:"port"`
}

type Configuration struct {
	Database DatabaseConfiguration `mapstructure:"database"`
	Server   ServerConfiguration   `mapstructure:"server"`
}

type Auth struct {
	AccessTokenExpirationDays  int    `mapstructure:"access_token_expiration_days"`
	RefreshTokenExpirationDays int    `mapstructure:"refresh_token_expiration_days"`
	secretKey                  string `mapstructure:"secret_key"`
}
