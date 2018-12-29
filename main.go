package main

import (
	"OpenStreetmapRouting/config"
	"OpenStreetmapRouting/data"
	"OpenStreetmapRouting/parsing"
	server "OpenStreetmapRouting/server"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {

	start := time.Now()
	config.LoadConfig("res/config.yaml")

	initLogger()

	InitGraphProd()

	log.Info("Ready!!")
	elapsed := time.Since(start)
	log.Infof("loading took %s", elapsed)

	server.Start()
}

func initLogger() {
	conf := config.GetConfig()

	var logLevel log.Level
	switch conf.LogLevel {
	case 1:
		logLevel = log.InfoLevel
	case 2:
		logLevel = log.WarnLevel
	case 3:
		logLevel = log.ErrorLevel
	default:
		logLevel = log.TraceLevel
	}

	log.SetReportCaller(true)

	log.SetLevel(logLevel)

}

//Init calculates the offsetlist and creates the grid for the given graph
func InitGraphProd() *data.GraphProd {

	conf := config.GetConfig()
	graphData := parsing.ParseOrLoadGraph(conf)

	g := data.InitGraphProd(graphData, conf)

	return g

}
