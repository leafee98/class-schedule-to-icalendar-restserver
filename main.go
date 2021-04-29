package main

import (
	"github.com/leafee98/class-schedule-to-icalendar-restserver/config"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/middlewares"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/routers"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/rpc"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/server"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error

	config.ParseParameter()
	if config.ConfigFile != "" {
		if err = config.LoadConfigFromFile(config.ConfigFile); err != nil {
			logrus.Fatal(err)
		}
	}

	config.LogCurrentConfig()

	if err = config.ValidParamCombination(); err != nil {
		logrus.Fatal(err)
	}

	// just initialize database or start server
	if config.InitDatabase {
		if err = db.InitDatabase(); err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("database initialized")
		return
	} else {
		if err = db.Init(); err != nil {
			logrus.Fatalf("failed to connect to database. detail: %s", err.Error())
		} else {
			logrus.Info("database connected")
		}
	}

	err = rpc.Init()
	if err != nil {
		logrus.Fatalf("Could not establish connection to %s. Please check the PRC server status. detail: %s",
			config.RPCTarget, err.Error())
	} else {
		logrus.Info("rpc server connected")
	}

	logrus.Info("starting rest server...")

	// keep this initialize order!
	server.Init()
	middlewares.Init(server.Engine)
	routers.Init(server.Engine)
	server.Run()
}
