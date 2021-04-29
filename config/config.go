package config

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// RPCTarget is the RPC server listen address and port, e: 127.0.0.1:8047
var RPCTarget string

// RestEndpoint is the endpoint of rest server's IP address, e: 0.0.0.0
var RestEndpoint string

// Database connect info
var DatabaseUsername string
var DatabasePassword string
var DatabaseHost string
var DatabaseName string

// HTTPBasepath is the base path while request this rest server, e: /api
var HTTPBasepath string

// if InitDatabase is set, just init database but don't start server
var InitDatabase bool

// specifiy the path of config file
var ConfigFile string

// const name of each configuration
type paramNames struct {
	DatabaseUsername string
	DatabasePassword string
	DatabaseHost     string
	DatabaseName     string

	RPCTarget    string
	RestEndpoint string
	HTTPBasepath string
	InitDatabase string
	ConfigFile   string
}

var pn paramNames = paramNames{
	DatabaseUsername: "database-username",
	DatabasePassword: "database-password",
	DatabaseHost:     "database-host",
	DatabaseName:     "database-name",

	RPCTarget:    "rpc-target",
	RestEndpoint: "rest-endpoint",
	HTTPBasepath: "http-basepath",
	InitDatabase: "init-database",
	ConfigFile:   "config",
}

// LoadConfig function load config from file whose path is confPath
func LoadConfigFromFile(confPath string) error {
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

				if err = loadSingleConfig(key, value); err != nil {
					return err
				}
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

func ParseParameter() {
	flag.StringVar(&DatabaseUsername, pn.DatabaseUsername, "", "username used to login database service.")
	flag.StringVar(&DatabasePassword, pn.DatabasePassword, "", "password used to login database service.")
	flag.StringVar(&DatabaseHost, pn.DatabaseHost, "", "host of database service. e: 127.0.0.1:3306")
	flag.StringVar(&DatabaseName, pn.DatabaseName, "", "name of database to use.")

	flag.StringVar(&RPCTarget, pn.RPCTarget, "", "RPCTarget is the RPC server listen address and port,"+
		"e: 127.0.0.1:8047")
	flag.StringVar(&RestEndpoint, pn.RestEndpoint, "", "RestEndpoint is the endpoint of rest server's IP address,"+
		"e: 0.0.0.0")
	flag.StringVar(&HTTPBasepath, pn.HTTPBasepath, "", "HTTPBasepath is the base path while request this rest server,"+
		"e: /api")
	flag.BoolVar(&InitDatabase, pn.InitDatabase, false, "add this parameter to init database and don't start server."+
		" (this parameter can noly specified in command line)")
	flag.StringVar(&ConfigFile, pn.ConfigFile, "", "specifiy the path of config file. "+
		" (this parameter can noly specified in command line)")
	flag.Parse()
}

func loadSingleConfig(key, value string) error {
	switch key {
	case pn.RPCTarget:
		if RPCTarget == "" {
			RPCTarget = value
		}
	case pn.RestEndpoint:
		if RestEndpoint == "" {
			RestEndpoint = value
		}
	case pn.HTTPBasepath:
		if HTTPBasepath == "" {
			HTTPBasepath = value
		}

	case pn.DatabaseUsername:
		if DatabaseUsername == "" {
			DatabaseUsername = value
		}
	case pn.DatabasePassword:
		if DatabasePassword == "" {
			DatabasePassword = value
		}
	case pn.DatabaseHost:
		if DatabaseHost == "" {
			DatabaseHost = value
		}
	case pn.DatabaseName:
		if DatabaseName == "" {
			DatabaseName = value
		}
	default:
		return errors.New(fmt.Sprintf("unrecognized: %s = %s", key, value))
	}
	return nil
}

func ValidParamCombination() error {
	if InitDatabase {
		if !validDatabaseSource() {
			return errors.New(fmt.Sprintf("when using %s, you must specify %s, %s, %s and %s",
				pn.InitDatabase, pn.DatabaseUsername, pn.DatabasePassword, pn.DatabaseHost, pn.DatabaseName))
		}
	} else {
		if !validDatabaseSource() ||
			HTTPBasepath == "" ||
			RestEndpoint == "" ||
			RPCTarget == "" {
			return errors.New("you haven't config all option")
		}
	}
	return nil
}

func validDatabaseSource() bool {
	return DatabaseUsername != "" &&
		DatabasePassword != "" &&
		DatabaseHost != "" &&
		DatabaseName != ""
}

func LogCurrentConfig() {
	logrus.Info("========== current config ==========")
	logrus.Infof("%20s = %s", pn.ConfigFile, ConfigFile)
	logrus.Infof("%20s = %t", pn.InitDatabase, InitDatabase)

	logrus.Infof("%20s = %s", pn.DatabaseUsername, DatabaseUsername)
	logrus.Infof("%20s = %s", pn.DatabasePassword, strings.Repeat("*", len(DatabasePassword)))
	logrus.Infof("%20s = %s", pn.DatabaseHost, DatabaseHost)
	logrus.Infof("%20s = %s", pn.DatabaseName, DatabaseName)

	logrus.Infof("%20s = %s", pn.RPCTarget, RPCTarget)
	logrus.Infof("%20s = %s", pn.RestEndpoint, RestEndpoint)
	logrus.Infof("%20s = %s", pn.HTTPBasepath, HTTPBasepath)
	logrus.Info("======== current config end =========")
}
