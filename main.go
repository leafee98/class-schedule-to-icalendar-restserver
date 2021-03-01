package main

import (
	"os"

	"github.com/leafee98/class-schedule-to-icalendar-restserver/config"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/middlewares"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/routers"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/rpc"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/server"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		logrus.Fatal("you must specify the config file path")
	}
	err := config.LoadConfig(os.Args[1])
	if err != nil {
		logrus.Fatalf("failed to load configure from file. detail: %s", err.Error())
	}

	err = db.Init()
	if err != nil {
		logrus.Fatalf("failed to connect to database. detail: %s", err.Error())
	} else {
		logrus.Info("database connected")
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
