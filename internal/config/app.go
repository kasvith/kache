package config

type AppConfig struct {
	Port       int
	Host       string
	Verbose    bool
	MaxClients int
	MaxTimeout int
}
