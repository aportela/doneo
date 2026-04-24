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
