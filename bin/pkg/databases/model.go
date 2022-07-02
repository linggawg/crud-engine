package databases

// DBConfig struct . .
type SQLXConfig struct {
	Host     string
	Port     uint16
	Name     string
	Username string
	Password string
	Dialect  string
	SSLMode  string
}
