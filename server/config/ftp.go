package config

type Ftp struct {
	Watch string `mapstructure:"watch" json:"watch" yaml:"watch"` // 监控目录
}
