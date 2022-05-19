package database

// DBConfig struct . .
type SQLXConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	Dialect  string
}

type MongoConfig struct {
	Host string
	Name string
}