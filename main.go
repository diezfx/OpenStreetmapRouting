package main

import (
	"fmt"
	"test/config"
	"test/data"
	"test/parsing"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {

	fmt.Println("Beginne")
	start := time.Now()
	config := config.LoadConfig("res/config.yaml")

	initLogger()

	graphData := parsing.ParseOrLoadGraph(config)

	//add the offset list that is needed for dijkstra
	graph := data.GraphProd{Nodes: graphData.Nodes, Edges: graphData.Edges}
	graph.CalcOffsetList()

	log.Info("Ready!!")
	elapsed := time.Since(start)
	log.Infof("loading took %s", elapsed)
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
