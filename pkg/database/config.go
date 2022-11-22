package database

type DatasourceConfig struct {
	DriverName string `mapstructure:"driverName"`
	Addr       string `mapstructure:"addr"`
	Database   string `mapstructure:"database"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	Charset    string `mapstructure:"charset"`
}
