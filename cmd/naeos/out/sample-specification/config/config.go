package config

type Config struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}
