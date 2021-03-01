package config

// RPCTarget is the RPC server listen address and port, e: 127.0.0.1:8047
var RPCTarget string

// RestEndpoint is the endpoint of rest server's IP address, e: 0.0.0.0
var RestEndpoint string

// DatabaseSource is the endpoint of mariadb/mysql, e: user:pass@127.0.0.1:3306/db
var DatabaseSource string

// HTTPBasepath is the base path while request this rest server, e: /api/
var HTTPBasepath string

// LoadConfig function load config from file whose path is confPath
// todo: load from file instead of using the predefine value
func LoadConfig(confPath string) error {
	RestEndpoint = "127.0.0.1:8048"
	RPCTarget = "127.0.0.1:8047"
	DatabaseSource = "u_csti:csti_pass@tcp(127.0.0.1:3306)/csti"
	HTTPBasepath = "/api"

	return nil
}
