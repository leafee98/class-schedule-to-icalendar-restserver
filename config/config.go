package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

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
	file, err := os.Open(confPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if (len(line) == 0 || line[0] == '#') && err == nil {
			continue
		}

		if symbolEqual := strings.Index(line, "="); symbolEqual >= 0 {
			if key := strings.TrimSpace(line[:symbolEqual]); len(key) > 0 {
				value := ""
				if len(line) > symbolEqual {
					value = strings.TrimSpace(line[symbolEqual+1:])
				}

				loadSingleConfig(key, value)
			}
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func loadSingleConfig(key, value string) error {
	switch key {
	case "rpc-target":
		RPCTarget = value
	case "rest-endpoint":
		RestEndpoint = value
	case "database-source":
		DatabaseSource = value
	case "http-basepath":
		HTTPBasepath = value
	default:
		return errors.New(fmt.Sprintf("unrecognized: %s = %s", key, value))
	}
	return nil
}
