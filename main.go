package main

import (
	"OpenStreetmapRouting/config"
	"OpenStreetmapRouting/controller"
	"OpenStreetmapRouting/data"
	"OpenStreetmapRouting/parsing"
	colorable "github.com/mattn/go-colorable"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {

	start := time.Now()
	conf := config.LoadConfig("res/config.yaml")

	initLogger()

	graph := InitGraphProd(conf)
	stations := data.InitGraphProdWithStations(graph, conf)

	// init grid for stations
	stationsGrid := InitStationsGrid(stations, conf)

	log.Info("Ready!!")
	elapsed := time.Since(start)
	log.Infof("loading took %s", elapsed)

	controller.Start(graph, stations, stationsGrid)
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
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())

	log.SetLevel(logLevel)

}

//InitGraphProd calculates the offsetlist and creates the grid for the given graph
func InitGraphProd(conf *config.Config) *data.GraphProd {

	graphData := parsing.ParseOrLoadGraph(conf)

	g := data.InitGraphProd(graphData, conf)

	return g

}

func InitStationsGrid(stations *data.GasStations, conf *config.Config) *data.Grid {
	stationsGrid := data.Grid{}
	stationList := make([]data.Node, len(stations.Stations))
	i := 0
	for _, value := range stations.Stations {
		stationList[i] = value
		i++
	}
	stationsGrid.InitGrid(stationList, conf)
	return &stationsGrid
}
