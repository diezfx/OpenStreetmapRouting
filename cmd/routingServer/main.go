package main

import (
	"github.com/diezfx/OpenStreetmapRouting/config"
	"github.com/diezfx/OpenStreetmapRouting/controller"
	"github.com/diezfx/OpenStreetmapRouting/data"
	"github.com/diezfx/OpenStreetmapRouting/parsing"
	"time"

	colorable "github.com/mattn/go-colorable"

	"github.com/sirupsen/logrus"
)

func main() {

	start := time.Now()
	conf := config.LoadConfig("res/config.yaml")

	initlogger()

	graph, stations := initGraphProd(conf)

	// init grid for stations
	stationsGrid := initStationsGrid(stations, conf)

	logrus.Info("Ready")
	elapsed := time.Since(start)
	logrus.Infof("loading took %s", elapsed)

	controller.Start(graph, stations, stationsGrid)
}

func initlogger() {
	conf := config.GetConfig()

	var logrusLevel logrus.Level
	switch conf.LogLevel {
	case 1:
		logrusLevel = logrus.InfoLevel
	case 2:
		logrusLevel = logrus.WarnLevel
	case 3:
		logrusLevel = logrus.ErrorLevel
	default:
		logrusLevel = logrus.TraceLevel
	}

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())

	logrus.SetLevel(logrusLevel)

}

//InitGraphProd calculates the offsetlist and creates the grid for the given graph
func initGraphProd(conf *config.Config) (*data.GraphProd, *data.GasStations) {

	graphData := parsing.ParseOrLoadGraph(conf)

	g := data.InitGraphProd(graphData, conf)
	grid := data.Grid{}
	grid.InitGrid(g.Nodes, conf)
	logrus.Debugf("Grid initialized")
	g.Grid = grid

	stations := data.InitGraphProdWithStations(g, conf)

	return g, stations

}

func initStationsGrid(stations *data.GasStations, conf *config.Config) *data.Grid {
	stationsGrid := data.Grid{}
	stationList := stations.GetStationsAsList()
	stationsGrid.InitGrid(stationList, conf)
	return &stationsGrid
}
